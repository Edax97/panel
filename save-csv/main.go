package main

import (
	"os"
)

func main() {
	var inputPath = os.Args[1]
	var savePath = os.Args[2]

	if savePath == "" {
		savePath = "csv-save"
	}
	if inputPath == "" {
		inputPath = "csv-input"
	}

	csvData := CSVData{saveComma: ',', inComma: ';'}
	GetEnergyCSV(csvData, inputPath, savePath)
}
