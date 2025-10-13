package main

import (
	"os"
)

func main() {
	var inputDir = os.Args[1]
	var saveDir = os.Args[2]

	if saveDir == "" {
		saveDir = "csv-save"
	}
	if inputDir == "" {
		inputDir = "csv-input"
	}

	dataSource := PowerData{saveComma: ',', inComma: ';'}
	FilterPower(dataSource, inputDir, saveDir)
}
