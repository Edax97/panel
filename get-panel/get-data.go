package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

var urlString = os.Getenv("PANEL_URLS")
var userString = os.Getenv("PANEL_USERS")
var passString = os.Getenv("PANEL_PASS")
var savePath = os.Getenv("CSV_INPUT_PATH")

func GetItems(list string) []string {
	return strings.Split(list, "\n")
}

type PanelStore interface {
	fetchCSV(url string, user string, pass string) (io.Reader, error)
	saveCSV(data io.Reader, fileName string) error
}

func GetPanelData(api PanelStore) {
	urlList := GetItems(urlString)
	userList := GetItems(userString)
	passList := GetItems(passString)

	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(len(urlList))

	for i, url := range urlList {
		go func() {
			data, err := api.fetchCSV(url, userList[i], passList[i])
			PanicError(err)
			mu.Lock()
			defer func() {
				mu.Unlock()
				wg.Done()
			}()

			err = api.saveCSV(data, fmt.Sprintf("%v/panel/%d.csvIO", savePath, i))
			PanicError(err)
			fmt.Println("Descarga exitosa de datos", url)
		}()
	}
	wg.Wait()
}
