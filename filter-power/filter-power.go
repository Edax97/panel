package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type IComServer interface {
	SendTimeValue(imei string, time time.Time, data string) (bool, error)
}
type IPanelStore interface {
	SendPanelServer(parsed [][]string, file string, serv IComServer) error
	SavePanelData(dir, file string)
}

func FilterPower(inputDir string, saveDir string, serv IComServer, p IPanelStore) error {

	var wg sync.WaitGroup

	files, err := os.ReadDir(inputDir)
	PanicError(err)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		wg.Add(1)
		go func(f string) {
			defer func() {
				wg.Done()
			}()
			if !strings.HasSuffix(f, ".csv") {
				fmt.Println("Ignoring file", f)
				return
			}
			data, err := ReadCSV(inputDir+"/"+f, ';')
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
	p.SavePanelData(saveDir, "panel_server.csv")
	return nil

}
