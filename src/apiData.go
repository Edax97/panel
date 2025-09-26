package main

import (
	filter_csv "data-store/filter-csv"
	"fmt"
	"strings"
	"sync"
)

// From env URLlist
const URLlist = "http://1.1.1.1\nhttp://1.2.3.1\nhttp://1.3.4.5"

//

func GetURL(list string) []string {
	return strings.Split(list, "\n")
}

func GetPanels() {
	csvStore := filter_csv.NewCSV(";", "storage")
	urlList := GetURL(URLlist)

	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(len(urlList))

	for _, url := range urlList {
		go func() {
			data, err := FetchDataCSV(url, ',')
			mu.Lock()
			defer mu.Unlock()
			// Filter data
			// Save data in file month, row day and column dev
			ok, err := csvStore.FilterSaveCSV(data)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(data, ok)
			wg.Done()
		}()
	}

	wg.Wait()
}
