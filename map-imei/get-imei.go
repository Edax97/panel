package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var imeiMap = make(map[string]string)
var imeiList = make([][]string, 1)

func main() {
	imeiList[0] = []string{"ID", "NAME", "PANEL", "IMEI"}
	InputDir := os.Args[1]
	ImeiFile := os.Args[2]

	var wg sync.WaitGroup
	var mutex sync.Mutex

	files, _ := os.ReadDir(InputDir)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		wg.Add(1)
		go func(f string) {
			fmt.Println(f)
			defer func() {
				mutex.Unlock()
				wg.Done()
			}()
			mutex.Lock()
			fileData, err := readCSV(InputDir+"/"+f, ',')
			if err != nil {
				log.Panic(err)
			}
			ids := fileData[0][1:]
			names := fileData[1][1:]

			for j, id := range ids {
				name := names[j]
				_, ok := imeiMap[id]
				if !ok {
					imeiList = append(imeiList, []string{id, name, f, ""})
					imeiMap[id] = name
				}
			}
		}(file.Name())
	}

	wg.Wait()
	print(len(imeiMap))
	saveSCV(ImeiFile, imeiList)
}
