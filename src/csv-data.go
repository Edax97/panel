package main

import (
	"encoding/csv"
	"fmt"
	"os"
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
	var month string
	var year, day int
	var WHCOLUMNS = make([]int, 0)

	deviceHeaders := parsed[2]
	fieldHeaders := parsed[4][1:]

	for i, field := range fieldHeaders {
		if field == "Total Received Active Energy" || field == "Total Delivered Active Energy" {
			WHCOLUMNS = append(WHCOLUMNS, i+1)
		}
	}

	for _, record := range parsed[6:] {
		timestamp := record[0]
		parsedTime, err := time.Parse("2006/01/02 15:04:05", timestamp)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			continue
		}
		if parsedTime.Minute() == 0 && parsedTime.Hour() == 0 {
			month = parsedTime.Month().String()
			year = parsedTime.Year()
			day = parsedTime.Day()

			saveFile := fmt.Sprintf("%s/%s-%d.csv", savePath, month, year)

			monthParsed, err := createMonthCSV(saveFile)
			if monthParsed == nil {
				monthParsed, err = ReadCSV(saveFile, C.saveComma)
			}
			PanicError(err)
			for _, index := range WHCOLUMNS {
				valueWh := record[index]
				device := deviceHeaders[index]
				// Copy valueWh to year-month file at row day and column device
				fmt.Printf("%v - Data: %s-%v-%v, %s: %s\n", index, month, year, day, device, valueWh)

				// Find device in file, otherwise add a column
				col := -1
				for i, header := range monthParsed[0] {
					if header == device {
						col = i
						break
					}
				}
				if col == -1 {
					col = len(monthParsed[0])
					monthParsed[0] = append(monthParsed[0], device)
					for j := 1; j < len(monthParsed); j++ {
						monthParsed[j] = append(monthParsed[j], "")
					}
				}
				monthParsed[day][col] = valueWh
			}

			SaveCSV(saveFile, monthParsed)
		}
	}

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
