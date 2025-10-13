package main

import (
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"time"
)

type SentCache struct {
	sentMap  map[string]time.Time
	diskPath string
}

func NewSentCache(diskPath string) *SentCache {
	cache := &SentCache{
		diskPath: diskPath,
		sentMap:  make(map[string]time.Time),
	}
	cache.loadCache()
	return cache
}

func (c *SentCache) saveCache() {
	f, err := os.Create(c.diskPath)
	if err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()
	encoder := gob.NewEncoder(f)
	err = encoder.Encode(c.sentMap)
	if err != nil {
		return
	}

}

func (c *SentCache) loadCache() {
	f, err := os.Open(c.diskPath)
	if err != nil {
		c.sentMap = make(map[string]time.Time)
		return
	}
	defer func() {
		_ = f.Close()
	}()

	var data map[string]time.Time
	decoder := gob.NewDecoder(f)
	err = decoder.Decode(&data)
	if err != nil {
		c.sentMap = make(map[string]time.Time)
		return
	}
	c.sentMap = data
}

func (c *SentCache) hasSent(imei string, t time.Time) bool {
	sent, ok := c.sentMap[imei]
	if !ok {
		return false
	}
	if t.Before(sent) {
		return true
	}
	return false
}

func (c *SentCache) updateSent(imei string, sent time.Time) bool {
	c.sentMap[imei] = sent
	c.saveCache()
	return true
}

type PowerData struct {
	inComma   rune
	saveComma rune
	server    ComServer
}

func (d PowerData) ReadCSVPower(fileName string) ([][]string, error) {
	return ReadCSV(fileName, d.inComma)
}

func (d PowerData) FilterPowerData(parsed [][]string, savePath string) error {
	//Sentcache
	cache := NewSentCache("sent-value.gob")

	var WHCOLUMNS = make([]int, 0)

	deviceHeaders := parsed[1]
	nameHeaders := parsed[2]
	fieldHeaders := parsed[4][1:]

	for i, field := range fieldHeaders {
		if field == "Total Received Active Energy" || field == "Total Delivered Active Energy" {
			WHCOLUMNS = append(WHCOLUMNS, i+1)
		}
	}
	filteredData := [][]string{{"ID"}, {"Time"}}
	for _, index := range WHCOLUMNS {
		device := deviceHeaders[index]
		name := nameHeaders[index]
		filteredData[0] = append(filteredData[0], device)
		filteredData[1] = append(filteredData[1], name)
	}

	for _, record := range parsed[6:] {
		timestamp := record[0]
		parsedTime, err := time.Parse("2006/01/02 15:04:05", timestamp)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			continue
		}
		if parsedTime.Minute() == 0 {
			row := []string{timestamp}
			for _, index := range WHCOLUMNS {
				valueWh := record[index]
				devId := deviceHeaders[index]
				imei := devId + "imei"

				intWh, err := strconv.Atoi(valueWh)
				if err != nil {
					fmt.Println("Error parsing value:", err)
					continue
				}
				row = append(row, valueWh)

				if cache.hasSent(imei, parsedTime) {
					continue
				}
				ok, _ := d.server.SendPowerValue(imei, parsedTime, intWh)
				if ok {
					cache.updateSent(imei, parsedTime)
				}
			}
			filteredData = append(filteredData, row)
		}
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
