package panelServer

import (
	"filter-power/csvIO"
	"filter-power/providers"
	"filter-power/wailonServer"
	"fmt"
	"log"
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
	imeiMap  map[string]*MedidorDatos
	transIds map[string]bool
	mutex    sync.Mutex
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
	return &PanelServer{imeiMap: imeiMap, transIds: transIds}, nil
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

				log.Printf("Id: %s %s", id, imei)
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
