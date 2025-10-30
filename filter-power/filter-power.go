package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type ComServer interface {
	SendTimeValue(imei string, time time.Time, wh string, vai string, vao string) (bool, error)
}

type CSVSource interface {
	ReadCSVPower(fileName string) ([][]string, error)
	SendWHData(data [][]string, dir string, file string) error
	SendHistoryWH(data [][]string, dir string, file string) error
}

func FilterPower(store CSVSource, inputDir string, saveDir string) {

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
			data, err := store.ReadCSVPower(inputDir + "/" + f)
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
			err = store.SendWHData(data, saveDir, f)
			//err = store.SendHistoryWH(data, saveDir, f)
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
		}(file.Name())
	}

	wg.Wait()

}
