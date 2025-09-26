package main

import (
	"fmt"
	"os"
	"time"
)

var DIR = os.Args[1]

var FILTERED = make([]int, 0)

func main() {

	if DIR == "" {
		DIR = "storage"
	}

	data, err := ReadDataDisc("data.csv", ';')
	if err != nil {
		panic(err)
		return
	}
	deviceHeaders := data[2]
	fieldHeaders := data[4][1:]

	for i, field := range fieldHeaders {
		if field == "Total Received Active Energy" || field == "Total Delivered Active Energy" {
			FILTERED = append(FILTERED, i+1)
		}
	}

	var month string
	var year, day int

	for _, record := range data[6:] {
		timestamp := record[0]
		parsedTime, err := time.Parse("2006/01/02 15:04:05", timestamp)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			continue
		}
		if parsedTime.Minute() == 0 && parsedTime.Hour() == 0 {
			month = parsedTime.Month().String()
			year = parsedTime.Year()
			day = parsedTime.Day()

			fileName := fmt.Sprintf("%s/%s-%d.csv", DIR, month, year)

			CreateMonthFile(fileName)

			//Open month file
			dataMonth, err := ReadDataDisc(fileName, ',')
			if err != nil {
				panic(err)
			}

			for _, index := range FILTERED {
				valueWh := record[index]
				device := deviceHeaders[index]
				// Copy valueWh to year-month file at row day and column device
				fmt.Printf("%v - Data: %s-%v-%v, %s: %s\n", index, month, year, day, device, valueWh)

				// Find device in file, otherwise add a column
				col := -1
				for i, header := range dataMonth[0] {
					if header == device {
						col = i
						break
					}
				}
				if col == -1 {
					col = len(dataMonth[0])
					dataMonth[0] = append(dataMonth[0], device)
					for j := 1; j < len(dataMonth); j++ {
						dataMonth[j] = append(dataMonth[j], "")
					}
				}
				dataMonth[day][col] = valueWh
			}

			WriteDataCSV(fileName, &dataMonth)
		}
	}
}
