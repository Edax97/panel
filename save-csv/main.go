package main

import (
	"os"
)

func main() {
	var inputPath = os.Getenv("CSV_INPUT_PATH")
	var savePath = os.Getenv("CSV_SAVE_PATH")

	if savePath == "" {
		savePath = "csv-save"
	}
	if inputPath == "" {
		inputPath = "csv-input"
	}

	csvData := CSVData{saveComma: ',', inComma: ';'}
	GetEnergyCSV(csvData, inputPath, savePath)
}
