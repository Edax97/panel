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

func (d PowerData) FilterPowerData(parsed [][]string, dir string, file string) error {
	// CACHE
	//cache := NewSentCache("sent-value.gob")
	savePath := fmt.Sprintf("%s/f_%s", dir, file)
	fmt.Println("Uploading data: ", file)
	if len(parsed) == 0 {
		return fmt.Errorf("empty file")
	}
	deviceHeaders := parsed[1]
	fieldHeaders := parsed[4][1:]
	devicePowerData := make(map[string]*struct {
		imei string
		j    int
		i    int
		data []string
	})

	for i, field := range fieldHeaders {
		// WHAT FIELD TO SEND
		if field == "Total Received Active Energy" || field == "Total Delivered Active Energy" {
			id := strings.Replace(deviceHeaders[i+1], "_O", "_I", 1)
			id = fmt.Sprintf("%s_%s", file, id)
			devData, ok := devicePowerData[id]
			if !ok {
				devicePowerData[id] = &struct {
					imei string
					j    int
					i    int
					data []string
				}{".", -1, -1, make([]string, 0)}
				devData = devicePowerData[id]
			}
			if field == "Total Received Active Energy" {
				devData.i = i + 1
			} else {
				devData.j = i + 1
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

		dat, ok := devicePowerData[id]
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
			for i, data := range devicePowerData {
				imeiParsed, err := strconv.Atoi(data.imei)
				if err != nil {
					//fmt.Println("Error parsing IMEI:", err)
					continue
				}
				imei := fmt.Sprintf("%d", 1e15+imeiParsed)[1:]

				v := getAtIndex(record, data.i)
				w := getAtIndex(record, data.j)
				if w > v {
					v = w
				}
				if v < 0 {
					v = 0
				}

				// CACHE
				//if cache.hasSent(imei, parsedTime) {
				//	continue
				//}
				//fmt.Printf("\n>>Sending to IMEI: %s | ID: %s | time %s\n",
				//	imei, id, timestamp)
				wg.Add(1)
				go func(IMEI string, WH int, id string) {
					defer wg.Done()
					ok, err := d.server.SendTimeValue(IMEI, parsedTime, WH)
					if !ok {
						log.Printf("Error sending: %s", err)
						return
					}
					mutex.Lock()
					defer mutex.Unlock()
					count++
					devicePowerData[i].data = append(devicePowerData[i].data, fmt.Sprintf("%s: %d", timestamp, WH))
				}(imei, v, i)
				// CACHE
				//	cache.updateSent(imei, parsedTime)
			}
			wg.Wait()
			fmt.Printf("> Panel %s | Time (%s) | Sent %d/%d\n", file, timestamp, count, len(devicePowerData))
		}
	}

	filteredData := [][]string{{"ID", "IMEI"}}
	for id, data := range devicePowerData {
		row := []string{id, data.imei}
		row = append(row, data.data...)
		filteredData = append(filteredData, row)
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

func getAtIndex(s []string, index int) int {
	if index < 0 || index >= len(s) {
		return -1
	}
	v, err := strconv.Atoi(s[index])
	if err != nil {
		return -1
	}
	return v
}
