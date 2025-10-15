package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type ComServer interface {
	SendTimeValue(imei string, time time.Time, value int) (bool, error)
}

type CSVSource interface {
	ReadCSVPower(fileName string) ([][]string, error)
	FilterPowerData(data [][]string, dir string, file string) error
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
			PanicError(err)
			err = store.FilterPowerData(data, saveDir, f)
			PanicError(err)
		}(file.Name())
	}

	wg.Wait()

}
