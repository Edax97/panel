package main

import (
	"os"
	"sync"
)

type CSVStore interface {
	ReadCSVInput(fileName string) ([][]string, error)
	FilterEnergyConsumption(data [][]string, savePath string) error
}

func GetEnergyCSV(store CSVStore, inputPath string, savePath string) {

	var wg sync.WaitGroup
	var mutex sync.Mutex

	files, err := os.ReadDir(inputPath)
	PanicError(err)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := inputPath + "/" + file.Name()
		wg.Add(1)
		go func(inputFile string) {
			defer func() {
				mutex.Unlock()
				wg.Done()
			}()
			mutex.Lock()
			data, err := store.ReadCSVInput(inputFile)
			PanicError(err)
			err = store.FilterEnergyConsumption(data, savePath)
			PanicError(err)
		}(fileName)
	}

	wg.Wait()

}
