package main

import (
	"filter-power/panelServer"
	"filter-power/wailonServer"
	"log"
	"os"
)

var IP string = os.Getenv("SERVER_IP")
var PORT string = os.Getenv("SERVER_PORT")

func main() {
	var inputDir = os.Args[1]
	var saveDir = os.Args[2]

	if saveDir == "" {
		saveDir = "../store/save"
	}
	if inputDir == "" {
		inputDir = "../store/input"
	}
	ser := wailonServer.NewWailonServer(IP, PORT)
	//ser := wailonServer.NewMockServer()
	p, err := panelServer.NewPanelServer()
	if err != nil {
		log.Fatalf("Error creating new panel server: %v", err)
	}
	err = FilterPower(inputDir, saveDir, ser, p)
	if err != nil {
		log.Fatalf("Error main loop: %v", err)
	}
}
