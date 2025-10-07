package main

import (
	"data-store/sendServer"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type CSVData struct {
	inComma   rune
	saveComma rune
}

func (C CSVData) ReadCSVInput(fileName string) ([][]string, error) {
	return ReadCSV(fileName, C.inComma)
}

func (C CSVData) FilterEnergyConsumption(parsed [][]string, savePath string) error {
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
			for j, index := range WHCOLUMNS {
				valueWh := record[index]
				intWh, err := strconv.Atoi(valueWh)
				if err != nil {
					fmt.Println("Error parsing value:", err)
					continue
				}

				row = append(row, valueWh)
				ok, _ := sendServer.SendMessage(deviceHeaders[j], parsedTime, intWh)
				if ok {
					//update lastSent value at imei = timestamp
				}
			}
			filteredData = append(filteredData, row)
		}
	}

	SaveCSV(savePath, filteredData)
	return nil
}

func createMonthCSV(name string) ([][]string, error) {
	f, _ := os.ReadFile(name)
	if len(f) > 0 {
		fmt.Println("File already exists")
		return nil, nil
	}

	fmt.Println("Creating file")
	file, err := os.Create(name)
	if err != nil {
		panic(err)
	}

	writer := csv.NewWriter(file)
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
		writer.Flush()
	}()

	dateCols := [][]string{[]string{"Día"}}
	for i := 1; i <= 31; i++ {
		date := []string{fmt.Sprintf("%d", i)}
		dateCols = append(dateCols, date)
	}
	err = writer.WriteAll(dateCols)
	if err != nil {
		panic(err)
	}
	return dateCols, nil
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
