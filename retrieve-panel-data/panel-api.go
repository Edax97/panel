package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
)

type PanelAPI struct {
	comma rune
}

func (p PanelAPI) fetchCSV(url string, user string, pass string) ([][]string, error) {
	//TODO Login
	//token := "ASDCASC1232c$35"
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
	reader.Comma = p.comma
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	return reader.ReadAll()
}

func (p PanelAPI) saveCSV(data [][]string, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)

	defer func() {
		_ = file.Close()
		writer.Flush()
	}()

	return writer.WriteAll(data)
}
