package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type ComServer interface {
	SendPowerValue(imei string, time time.Time, value int) (bool, error)
}

type CSVSource interface {
	ReadCSVPower(fileName string) ([][]string, error)
	FilterPowerData(data [][]string, savePath string) error
}

func FilterPower(store CSVSource, inputDir string, saveDir string) {

	var wg sync.WaitGroup
	var mutex sync.Mutex

	files, err := os.ReadDir(inputDir)
	PanicError(err)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		wg.Add(1)
		go func(f string) {
			defer func() {
				mutex.Unlock()
				wg.Done()
			}()
			mutex.Lock()
			data, err := store.ReadCSVPower(inputDir + "/" + f)
			PanicError(err)
			savePath := fmt.Sprintf("%s/f_%s", saveDir, f)
			err = store.FilterPowerData(data, savePath)
			PanicError(err)
		}(file.Name())
	}

	wg.Wait()

}
