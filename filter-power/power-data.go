package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MedidorDatos struct {
	imei    string
	colWH   int
	colVAI  int
	colVAO  int
	dataWH  []string
	dataVAI []string
	dataVAO []string
}

func NewMedidorDatos() *MedidorDatos {
	return &MedidorDatos{".", -1, -1, -1,
		[]string{}, []string{}, []string{}}
}

type PowerData struct {
	inComma   rune
	saveComma rune
	server    ComServer
}

func NewPowerData(server ComServer) PowerData {
	return PowerData{';', ',', server}
}

func (d PowerData) ReadCSVPower(fileName string) ([][]string, error) {
	return ReadCSV(fileName, d.inComma)
}

func (d PowerData) SendWHData(parsed [][]string, dir string, file string) error {
	// CACHE
	//cache := NewSentCache("sent-value.gob")
	savePath := fmt.Sprintf("%s/f_%s", dir, file)
	fmt.Println("Uploading Imei: ", file)
	if len(parsed) == 0 {
		return fmt.Errorf("empty file")
	}
	// Id -> medidor index
	deviceHeaders := parsed[1]
	fieldHeaders := parsed[4][1:]
	idToMedidor := make(map[string]*MedidorDatos)
	for i, field := range fieldHeaders {
		// WHAT FIELD TO SEND
		if field == "Total Delivered Active Energy" ||
			field == "Total Delivered Reactive Energy" ||
			field == "Total Received Reactive Energy" {
			idParts := strings.Split(deviceHeaders[i+1], "_")
			if len(idParts) < 2 {
				continue
			}
			id := fmt.Sprintf("%s_%s_%s_WHr_I", file, idParts[0], idParts[1])
			vals, ok := idToMedidor[id]
			if !ok {
				idToMedidor[id] = NewMedidorDatos()
				vals = idToMedidor[id]
			}
			if field == "Total Delivered Active Energy" {
				vals.colWH = i + 1
			} else if field == "Total Delivered Reactive Energy" {
				vals.colVAI = i + 1
			} else {
				vals.colVAO = i + 1
			}
		}
	}

	// Id -> imei
	imeiMap := os.Getenv("IMEI_MAP")
	imeiList := strings.Split(imeiMap, "\n")
	if len(imeiList) == 0 {
		return fmt.Errorf("imei file not set")
	}
	for _, line := range imeiList {
		v := strings.Split(line, ",")
		id, imei := v[0], v[2]
		dat, ok := idToMedidor[id]
		if ok {
			dat.imei = imei
		}
	}

	for _, record := range parsed[6:] {
		timestamp := record[0]
		loc, _ := time.LoadLocation("America/Lima")
		parsedTime, err := time.ParseInLocation("2006/01/02 15:04:05", timestamp, loc)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			continue
		}
		if parsedTime.Minute() == 0 {
			count := 0
			var wg sync.WaitGroup
			var mutex sync.Mutex
			for id, devInfo := range idToMedidor {
				imeiParsed, err := strconv.Atoi(devInfo.imei)
				if err != nil {
					//fmt.Println("Error parsing IMEI:", err)
					continue
				}
				imei := fmt.Sprintf("%d", 1e15+imeiParsed)[1:]

				wh := getValueAt(record, devInfo.colWH)
				vai := getValueAt(record, devInfo.colVAI)
				vao := getValueAt(record, devInfo.colVAO)

				// CACHE
				//if cache.hasSent(imei, parsedTime) {
				//	continue
				//}
				//fmt.Printf("\n>>Sending to IMEI: %s | ID: %s | time %s\n",
				//	imei, id, timestamp)
				wg.Add(1)
				go func(IMEI string, ID string, WH string, VAI string, VAO string) {
					defer wg.Done()
					ok, err := d.server.SendTimeValue(IMEI, parsedTime, WH, VAI, VAO)
					if !ok {
						log.Printf("Error sending: %s", err)
						return
					}
					mutex.Lock()
					defer mutex.Unlock()
					count++
					idToMedidor[ID].dataWH = append(idToMedidor[ID].dataWH, fmt.Sprintf("%s: %s", timestamp, WH))
					idToMedidor[ID].dataVAI = append(idToMedidor[ID].dataVAI, fmt.Sprintf("%s: %s", timestamp, VAI))
					idToMedidor[ID].dataVAO = append(idToMedidor[ID].dataVAO, fmt.Sprintf("%s: %s", timestamp, VAO))
				}(imei, id, wh, vai, vao)
				// CACHE
				//	cache.updateSent(imei, parsedTime)
			}
			wg.Wait()
			fmt.Printf("> Panel %s | Time (%s) | Sent %d/%d\n", file, timestamp, count, len(idToMedidor))
		}
	}

	filteredData := [][]string{{"ID", "IMEI", "VARIABLE"}}
	for id, Imei := range idToMedidor {
		rowh := []string{id, Imei.imei, "WH"}
		rowh = append(rowh, Imei.dataWH...)
		rvai := []string{id, Imei.imei, "VARH I"}
		rvai = append(rvai, Imei.dataVAI...)
		rvao := []string{id, Imei.imei, "VARH O"}
		rvao = append(rvao, Imei.dataVAO...)
		filteredData = append(filteredData, rowh, rvai, rvao)
	}

	SaveCSV(savePath, filteredData)
	return nil
}

func (d PowerData) SendHistoryWH(parsed [][]string, dir string, file string) error {
	// CACHE
	//cache := NewSentCache("sent-value.gob")
	savePath := fmt.Sprintf("%s/f_%s", dir, file)
	fmt.Println("Uploading data: ", file)
	if len(parsed) == 0 {
		return fmt.Errorf("empty file")
	}
	imeiHeaders := parsed[3][1:]
	locHeaders := parsed[4][1:]

	devDataMap := make(map[string]*struct {
		imei    string
		dataWH  []string
		dataVAH []string
	})
	for j, l := range locHeaders {
		devDataMap[l] = &struct {
			imei    string
			dataWH  []string
			dataVAH []string
		}{
			imeiHeaders[j],
			[]string{},
			[]string{},
		}
	}

	for _, record := range parsed[6:] {
		timestamp := record[0]
		loc, _ := time.LoadLocation("America/Lima")
		parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", timestamp, loc)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			continue
		}
		if parsedTime.Minute() == 0 {
			//fmt.Println("\nAt ", timestamp)
			count := 0
			// Concurrently send devs
			var wg sync.WaitGroup
			var mutex sync.Mutex
			for i, imeiNum := range imeiHeaders {
				loc := locHeaders[i]

				wh := getValueAt(record, i+1)

				imeiParsed, err := strconv.Atoi(imeiNum)
				if err != nil {
					//fmt.Println("Error parsing IMEI:", err)
					continue
				}
				imei := fmt.Sprintf("%d", 1e15+imeiParsed)[1:]

				wg.Add(1)
				go func(IMEI string, WH string, LOC string) {
					defer wg.Done()
					ok, err := d.server.SendTimeValue(IMEI, parsedTime, WH, "NaN", "NaN")
					if !ok {
						log.Printf("Error sending: %s", err)
						return
					}
					mutex.Lock()
					defer mutex.Unlock()

					devDataMap[loc].dataWH = append(devDataMap[loc].dataWH, fmt.Sprintf("%s: %s", timestamp, WH))
					devDataMap[loc].dataVAH = append(devDataMap[loc].dataVAH, fmt.Sprintf("%s: %s", timestamp, "NaN"))

					count++
				}(imei, wh, loc)
				// CACHE
				//	cache.updateSent(imei, parsedTime)
			}
			wg.Wait()
			fmt.Printf("> Panel %s | Time (%s) | Sent %d\n", file, timestamp, count)
		}
	}

	filteredData := [][]string{{"LOC", "IMEI", "VARIABLE"}}
	for loc, dev := range devDataMap {
		rowh := []string{loc, dev.imei, "WH"}
		rowh = append(rowh, dev.dataWH...)
		rvarh := []string{loc, dev.imei, "VARH"}
		rvarh = append(rvarh, dev.dataVAH...)
		filteredData = append(filteredData, rowh, rvarh)
	}
	SaveCSV(savePath, filteredData)
	return nil
}
func ReadCSV(name string, comma rune) ([][]string, error) {
	file, err := os.Open(name)
	if err != nil {
		fmt.Printf("file not found: %s\n", err)
		return nil, err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()
	reader := csv.NewReader(file)
	reader.Comma = comma
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	return reader.ReadAll()
}

func SaveCSV(name string, data [][]string) {
	file, err := os.Create(name)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	writer := csv.NewWriter(file)

	defer func() {
		_ = file.Close()
		writer.Flush()
	}()
	err = writer.WriteAll(data)
	if err != nil {
		panic(err)
	}
}

func getValueAt(s []string, index int) string {
	if index < 0 || index >= len(s) {
		return "NaN"
	}
	v, err := strconv.Atoi(s[index])
	if err != nil {
		return "NaN"
	}
	return fmt.Sprintf("%d", v)
}
