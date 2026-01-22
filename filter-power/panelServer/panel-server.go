package panelServer

import (
	"errors"
	"filter-power/csvIO"
	"filter-power/providers"
	"filter-power/wailonServer"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MedidorDatos struct {
	imei string
	data []string
}
type DeviceFieldNames struct {
	cols  []string
	mutex sync.Mutex
}

type TransformerDatum struct {
	params TransformerParams
	time   time.Time
}

func (f *DeviceFieldNames) getField(record []string, id, name string) string {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	//id := fmt.Sprintf("%s_%s_%s_WHr_I", file, idParts[0], idParts[1])
	parts := strings.Split(id, "_")
	dev := parts[2] + "_" + parts[3]
	col := dev + "_" + name

	index := -1
	for i, c := range f.cols {
		if c == col {
			index = i
			break
		}
	}
	if index < 0 {
		return "NaN"
	}
	v, err := strconv.ParseFloat(record[index], 64)
	if err != nil {
		return "NaN"
	}
	return fmt.Sprintf("%f", v)
}

type PanelServer struct {
	imeiMap        map[string]*MedidorDatos
	transIds       map[string]bool
	transformerMap map[string][]TransformerDatum
	mutex          sync.Mutex
}

func NewPanelServer() (*PanelServer, error) {
	// Id -> imei
	imeiMap := make(map[string]*MedidorDatos)
	imeiFile := os.Getenv("IMEI_MAP")
	imeiList := strings.Split(imeiFile, "\n")
	if len(imeiList) == 0 {
		return nil, fmt.Errorf("imei file not set")
	}
	for _, line := range imeiList {
		v := strings.Split(line, ",")
		id, imei := v[0], v[2]
		_, ok := imeiMap[id]
		if !ok {
			imeiMap[id] = &MedidorDatos{imei, []string{}}
		}
	}
	// ids of transformers
	transformerIdStr := os.Getenv("TRANSFORMER_IDS")
	transIds := make(map[string]bool)
	for _, id := range strings.Split(transformerIdStr, "\n") {
		transIds[id] = true
	}
	transformerMap := make(map[string][]TransformerDatum, 2)
	return &PanelServer{imeiMap: imeiMap, transIds: transIds, transformerMap: transformerMap}, nil
}

func (p *PanelServer) IsTransformer(id string) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	_, ok := p.transIds[id]
	return ok
}
func (p *PanelServer) getIdToImei(devHeads []string, file string) (map[string]string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	imeis := make(map[string]string)
	for _, dev := range devHeads {
		idParts := strings.Split(dev, "_")
		if len(idParts) < 2 {
			continue
		}
		id := fmt.Sprintf("%s_%s_%s_WHr_I", file, idParts[0], idParts[1])
		_, ok := imeis[id]
		if !ok {
			dev, ok := p.imeiMap[id]
			if !ok {
				continue
			}
			imeiParsed, err := strconv.Atoi(dev.imei)
			if err != nil {
				continue
			}
			imeis[id] = fmt.Sprintf("%d", 1e15+imeiParsed)[1:]
		}
	}
	return imeis, nil
}

func (p *PanelServer) SendPanelServer(parsed [][]string, file string, serv providers.IComServer) error {
	// CACHE
	cache := wailonServer.NewSentCache(file + "_cache.gob")
	fmt.Println("Uploading Imei: ", file)
	if len(parsed) == 0 {
		return fmt.Errorf("empty file")
	}

	deviceHeaders := parsed[1]
	devFields := DeviceFieldNames{cols: deviceHeaders}
	idToImei, err := p.getIdToImei(deviceHeaders, file)
	if err != nil {
		return err
	}

	for _, record := range parsed[6:] {
		timestamp := record[0]
		loc, _ := time.LoadLocation("America/Lima")
		parsedTime, err := time.ParseInLocation("2006/01/02 15:04:05", timestamp, loc)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			continue
		}
		// CACHE
		if cache.HasSent(file, parsedTime) {
			continue
		}
		if parsedTime.Minute() == 0 {
			count := 0
			var wg sync.WaitGroup
			for id, imei := range idToImei {

				//log.Printf("Id: %s %s", id, imei)
				wg.Add(1)
				go func(IMEI string, ID string, row []string) {
					defer wg.Done()
					var data string

					if !p.IsTransformer(ID) {
						wh := devFields.getField(row, ID, "WHr_I")
						vai := devFields.getField(row, ID, "VARHr_I")
						vao := devFields.getField(row, ID, "VARHr_O")
						data = fmt.Sprintf("watth:3:%s,varh:3:%s,varo:3:%s;", wh, vai, vao)
					} else {
						ia := devFields.getField(row, ID, "IA")
						ib := devFields.getField(row, ID, "IB")
						ic := devFields.getField(row, ID, "IC")
						vab := devFields.getField(row, ID, "VAB")
						vbc := devFields.getField(row, ID, "VBC")
						vca := devFields.getField(row, ID, "VCA")
						van := devFields.getField(row, ID, "VAN")
						vbn := devFields.getField(row, ID, "VBN")
						vcn := devFields.getField(row, ID, "VCN")
						wh := devFields.getField(row, ID, "WHr_I")
						vai := devFields.getField(row, ID, "VARHr_I")
						vao := devFields.getField(row, ID, "VARHr_O")
						vain := devFields.getField(row, ID, "VAHrIn")
						pftl := devFields.getField(row, ID, "PFTtl")

						transformer := CalcTransformer(ia, ib, ic, vab, vbc, vca, pftl)
						if _, ok := p.transformerMap[ID]; ok {
							p.transformerMap[ID] = append(p.transformerMap[ID],
								TransformerDatum{
									time:   parsedTime,
									params: *transformer,
								})
						} else {
							p.transformerMap[ID] = []TransformerDatum{
								{
									time:   parsedTime,
									params: *transformer,
								},
							}
						}

						dataSec := fmt.Sprintf("watth:3:%s,varh:3:%s,varo:3:%s,", wh, vai, vao) +
							fmt.Sprintf("varin:3:%s,pfttl:3:%s,", vain, pftl) +
							fmt.Sprintf("Ia:3:%s,Ib:3:%s,Ic:3:%s,", ia, ib, ic) +
							fmt.Sprintf("Vab:3:%s,Vbc:3:%s,Vca:3:%s,", vab, vbc, vca) +
							fmt.Sprintf("Van:3:%s,Vbn:3:%s,Vcn:3:%s", van, vbn, vcn)
						dataPrim := fmt.Sprintf("Iaprim:3:%s,Ibprim:3:%s,Icprim:3:%s,Vabprim:3:%s,Vbcprim:3:%s,Vcaprim:3:%s,Pprim:3:%s,Qprim:3:%s,Sprim:3:%s",
							transformer.IaPrim, transformer.IbPrim, transformer.IcPrim,
							transformer.VabPrim, transformer.VbcPrim, transformer.VcaPrim,
							transformer.Pprim, transformer.Qprim, transformer.Sprim,
						)
						data = dataSec + "," + dataPrim + ";"
					}

					ok, err := serv.SendTimeValue(IMEI, parsedTime, data)
					if !ok {
						log.Printf("Error sending: %s", err)
						return
					}
					p.mutex.Lock()
					defer p.mutex.Unlock()
					count++
					p.imeiMap[ID].data = append(p.imeiMap[ID].data, fmt.Sprintf("%s: %s", timestamp, data))
				}(imei, id, record)
			}
			wg.Wait()

			// CACHE
			cache.UpdateSent(file, parsedTime)
			fmt.Printf("> Panel %s | Time (%s) | Sent %d/%d\n", file, timestamp, count, len(idToImei))
		}
	}

	return nil
}

func (p *PanelServer) SavePanelData(dir, file string) {
	filteredData := [][]string{{"ID", "IMEI", "DATA"}}
	for id, Imei := range p.imeiMap {
		rowh := []string{id, Imei.imei}
		rowh = append(rowh, Imei.data...)
		filteredData = append(filteredData, rowh)
	}
	csvIO.SaveCSV(fmt.Sprintf("%s/%s", dir, file), filteredData)
}

func (p *PanelServer) SendCeldaRemonte(serv providers.IComServer) error {
	imei := os.Getenv("CELDA_REMONTE_IMEI")
	if imei == "" {
		return fmt.Errorf("CELDA_REMONTE_IMEI not loaded")
	}

	var TransformerIds = os.Getenv("TRANSFORMER_IDS")
	idList := strings.Split(TransformerIds, "\n")
	if len(idList) != 2 {
		return fmt.Errorf("not enough ids: %v", idList)
	}

	valuesT1, ok := p.transformerMap[idList[1]]
	if !ok {
		return fmt.Errorf("transformer 1 (%s) values not mapped", idList[1])
	}
	valuesT2, ok := p.transformerMap[idList[0]]
	if !ok {
		return fmt.Errorf("transformer 2 (%s) values not mapped", idList[0])
	}
	if len(valuesT1) != len(valuesT2) {
		return errors.New("different measurements")
	}

	for j, datum1 := range valuesT1 {
		datum2 := valuesT2[j]

		t1 := datum1.params
		t2 := datum2.params

		Iacr := toFloat(t1.IaPrim) + toFloat(t2.IaPrim) //IaPrim t1 +IaPrim t2
		Ibcr := toFloat(t1.IbPrim) + toFloat(t2.IbPrim) //IbPrim t1 + IbPrim t2
		Iccr := toFloat(t1.IcPrim) + toFloat(t2.IcPrim) //IcPrim t1 + IcPrim t2
		Vabcr := math.Max(toFloat(t1.VabPrim), toFloat(t2.VabPrim))
		Vbccr := math.Max(toFloat(t1.VbcPrim), toFloat(t2.VbcPrim))
		Vcacr := math.Max(toFloat(t1.VcaPrim), toFloat(t2.VcaPrim))
		Iavgcr := (Iacr + Ibcr + Iccr) / 3
		Vavgcr := (Vabcr + Vbccr + Vcacr) / 3
		Scr := Iavgcr * Vavgcr * math.Sqrt(3)
		Pcr := toFloat(t1.Pprim) + toFloat(t2.Pprim)
		Qcr := toFloat(t1.Qprim) + toFloat(t2.Qprim)

		data := fmt.Sprintf("iacr:3:%.2f,ibcr:3:%.2f,iccr:3:%.2f,vabcr:3:%.2f,vbccr:3:%.2f,vcacr:3:%.2f,iavgcr:3:%.2f,vavgcr:3:%.2f,scr:3:%.2f,pcr:3:%.2f,qcr:3:%.2f;",
			Iacr, Ibcr, Iccr,
			Vabcr, Vbccr, Vcacr, Iavgcr, Vavgcr,
			Scr, Pcr, Qcr)

		if ok, err := serv.SendTimeValue(imei, datum1.time, data); !ok {
			log.Printf("error sending celda remonte: %v", err)
		}
	}

	return nil
}
func toFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
