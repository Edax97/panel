package main

import (
	"filter-power/wailonServer"
	"os"
)

var IP string = os.Getenv("SERVER_IP")
var PORT string = os.Getenv("SERVER_PORT")

func main() {
	var inputDir = os.Args[1]
	var saveDir = os.Args[2]

	if saveDir == "" {
		saveDir = "csv-save"
	}
	if inputDir == "" {
		inputDir = "csv-input"
	}
	//ser := wailonServer.NewWailonServer(IP, PORT)
	ser := wailonServer.NewMockServer()
	d := NewPowerData(ser)

	FilterPower(d, inputDir, saveDir)
}
