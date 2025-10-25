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
	deviceHeaders := parsed[1]
	fieldHeaders := parsed[4][1:]
	idToValues := make(map[string]*struct {
		imei     string
		colWH    int
		colVAR   int
		dataWH   []string
		dataVARH []string
	})
	for i, field := range fieldHeaders {
		// WHAT FIELD TO SEND
		if field == "Total Delivered Active Energy" || field == "Total Delivered Reactive Energy" {
			id := strings.Replace(deviceHeaders[i+1], "_O", "_I", 1)
			id = fmt.Sprintf("%s_%s", file, id)
			vals, ok := idToValues[id]
			if !ok {
				idToValues[id] = &struct {
					imei     string
					colWH    int
					colVAR   int
					dataWH   []string
					dataVARH []string
				}{".", -1, -1, []string{}, []string{}}
				vals = idToValues[id]
			}
			if field == "Total Delivered Active Energy" {
				vals.colWH = i + 1
			} else {
				vals.colVAR = i + 1
			}
		}
	}

	// Set IMEI at device id
	imeiMap := os.Getenv("IMEI_MAP")
	//imeiMap := IMEI_MAP
	imeiList := strings.Split(imeiMap, "\n")

	//fmt.Println("Line", imeiMap)
	if len(imeiList) == 0 {
		return fmt.Errorf("imei file not set")
	}
	for _, line := range imeiList {
		v := strings.Split(line, ",")
		id, imei := v[0], v[2]
		dat, ok := idToValues[id]
		if ok {
			dat.imei = imei
		}
	}

	for _, record := range parsed[6:] {
		timestamp := record[0]
		parsedTime, err := time.Parse("2006/01/02 15:04:05", timestamp)
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
			for id, Imei := range idToValues {
				imeiParsed, err := strconv.Atoi(Imei.imei)
				if err != nil {
					//fmt.Println("Error parsing IMEI:", err)
					continue
				}
				imei := fmt.Sprintf("%d", 1e15+imeiParsed)[1:]

				wh := getValueAt(record, Imei.colWH)
				varh := getValueAt(record, Imei.colVAR)

				// CACHE
				//if cache.hasSent(imei, parsedTime) {
				//	continue
				//}
				//fmt.Printf("\n>>Sending to IMEI: %s | ID: %s | time %s\n",
				//	imei, id, timestamp)
				wg.Add(1)
				go func(IMEI string, WH string, VARH string, ID string) {
					defer wg.Done()
					ok, err := d.server.SendTimeValue(IMEI, parsedTime, WH, VARH)
					if !ok {
						log.Printf("Error sending: %s", err)
						return
					}
					mutex.Lock()
					defer mutex.Unlock()
					count++
					idToValues[ID].dataWH = append(idToValues[ID].dataWH, fmt.Sprintf("%s: %s", timestamp, WH))
					idToValues[ID].dataVARH = append(idToValues[ID].dataVARH, fmt.Sprintf("%s: %s", timestamp, VARH))
				}(imei, wh, varh, id)
				// CACHE
				//	cache.updateSent(imei, parsedTime)
			}
			wg.Wait()
			fmt.Printf("> Panel %s | Time (%s) | Sent %d/%d\n", file, timestamp, count, len(idToValues))
		}
	}

	filteredData := [][]string{{"ID", "IMEI", "VARIABLE"}}
	for id, Imei := range idToValues {
		rowh := []string{id, Imei.imei, "WH"}
		rowh = append(rowh, Imei.dataWH...)
		rvarh := []string{id, Imei.imei, "VARH"}
		rvarh = append(rvarh, Imei.dataVARH...)
		filteredData = append(filteredData, rowh, rvarh)
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
		parsedTime, err := time.Parse("2006-01-02 15:04:05", timestamp)
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
					ok, err := d.server.SendTimeValue(IMEI, parsedTime, WH, "NaN")
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
