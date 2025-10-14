package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

var imeiMap = make(map[string]string)
var imeiList = make([][]string, 1)

func main() {
	imeiList[0] = []string{"ID", "NAME", "IMEI"}
	InputDir := os.Args[1]
	ImeiFile := os.Args[2]

	var wg sync.WaitGroup
	var mutex sync.Mutex

	files, _ := os.ReadDir(InputDir)

	fmt.Println(len(files))
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
			fileData, err := readCSV(InputDir+"/"+f, ';')
			if err != nil {
				log.Panic(err)
			}
			devs := fileData[1][1:]
			names := fileData[2][1:]

			for j, d := range devs {
				name := names[j]
				if !strings.HasSuffix(d, "WHr_I") && !strings.HasSuffix(d, "WHr_O") {
					continue
				}
				id := fmt.Sprintf("%s_%s", f, d)
				_, ok := imeiMap[id]
				if !ok {
					imeiList = append(imeiList, []string{id, name, ""})
					imeiMap[id] = name
				}
			}
		}(file.Name())
	}

	wg.Wait()
	print(len(imeiMap))
	saveSCV(ImeiFile, imeiList)
}
