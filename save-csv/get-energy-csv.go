package main

import (
	"fmt"
	"os"
	"sync"
)

type CSVStore interface {
	ReadCSVInput(fileName string) ([][]string, error)
	FilterEnergyConsumption(data [][]string, savePath string) error
}

func GetEnergyCSV(store CSVStore, inputDir string, saveDir string) {

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
			data, err := store.ReadCSVInput(inputDir + "/" + f)
			PanicError(err)
			savePath := fmt.Sprintf("%s/filtered_%s", saveDir, f)
			err = store.FilterEnergyConsumption(data, savePath)
			PanicError(err)
		}(file.Name())
	}

	wg.Wait()

}
