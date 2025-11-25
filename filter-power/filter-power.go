package main

import (
	"filter-power/csvIO"
	"filter-power/providers"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func FilterPower(inputDir string, saveDir string, serv providers.IComServer, p providers.IPanelStore) error {

	var wg sync.WaitGroup

	files, err := os.ReadDir(inputDir)
	if err != nil {
		log.Panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		wg.Add(1)
		go func(f string) {
			defer func() {
				wg.Done()
			}()
			if !strings.HasSuffix(f, ".csvIO") {
				fmt.Println("Ignoring file", f)
				return
			}
			data, err := csvIO.ReadCSV(inputDir+"/"+f, ';')
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
			err = p.SendPanelServer(data, f, serv)
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
		}(file.Name())
	}

	wg.Wait()
	p.SavePanelData(saveDir, "panel_server.csvIO")
	return nil

}
