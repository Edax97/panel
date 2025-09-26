package main

import (
	"os"
)

func main() {
	var inputPath = os.Getenv("CSV_INPUT_PATH")
	var savePath = os.Getenv("CSV_SAVE_PATH")

	if savePath == "" {
		savePath = "save-csv"
	}
	if inputPath == "" {
		inputPath = "input-csv"
	}

	csvData := CSVData{saveComma: ',', inComma: ';'}
	GetEnergyCSV(csvData, inputPath, savePath)
}
