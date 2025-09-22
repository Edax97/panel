package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func CreateMonthFile(name string) {
	f, _ := os.ReadFile(name)
	if len(f) > 0 {
		fmt.Println("File already exists")
		return
	}

	fmt.Println("Creating file")
	file, err := os.Create(name)
	if err != nil {
		panic(err)
		return
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
}

func ReadDataDisc(name string, comma rune) ([][]string, error) {
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

func WriteDataCSV(name string, data *[][]string) {
	// Write data to file
	file, err := os.Create(name)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	writer := csv.NewWriter(file)

	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
		writer.Flush()
	}()

	err = writer.WriteAll(*data)
	if err != nil {
		panic(err)
	}
}
