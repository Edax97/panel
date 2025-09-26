package filter_csv

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
)

func CreateMonthCSV(name string) ([][]string, error) {
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

func FetchDataCSV(url string, comma rune) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)
	reader.Comma = comma
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	return reader.ReadAll()
}

func SaveCSV(name string, data *[][]string) {
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
