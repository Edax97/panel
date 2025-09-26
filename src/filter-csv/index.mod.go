package filter_csv

import (
	"fmt"
	"time"
)

type CSVFiles struct {
	comma string
	dir   string
}

type Filter interface {
	FilterSaveCSV(parsed [][]string) (bool, error)
}

func NewCSV(comma string, dir string) *CSVFiles {
	if dir == "" {
		dir = "storage"
	}
	return &CSVFiles{comma, dir}
}

func (c *CSVFiles) FilterSaveCSV(parsed [][]string) (bool, error) {
	var month string
	var year, day int
	var WHCOLUMNS = make([]int, 0)

	deviceHeaders := parsed[2]
	fieldHeaders := parsed[4][1:]

	for i, field := range fieldHeaders {
		if field == "Total Received Active Energy" || field == "Total Delivered Active Energy" {
			WHCOLUMNS = append(WHCOLUMNS, i+1)
		}
	}

	for _, record := range parsed[6:] {
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

			monthParsed, err := CreateMonthCSV(fileName)
			if monthParsed == nil {
				monthParsed, err = ReadCSV(fileName, ',')
			}
			if err != nil {
				panic(err)
			}

			for _, index := range WHCOLUMNS {
				valueWh := record[index]
				device := deviceHeaders[index]
				// Copy valueWh to year-month file at row day and column device
				fmt.Printf("%v - Data: %s-%v-%v, %s: %s\n", index, month, year, day, device, valueWh)

				// Find device in file, otherwise add a column
				col := -1
				for i, header := range monthParsed[0] {
					if header == device {
						col = i
						break
					}
				}
				if col == -1 {
					col = len(monthParsed[0])
					monthParsed[0] = append(monthParsed[0], device)
					for j := 1; j < len(monthParsed); j++ {
						monthParsed[j] = append(monthParsed[j], "")
					}
				}
				monthParsed[day][col] = valueWh
			}

			SaveCSV(fileName, &monthParsed)
		}
	}

	return true, nil
}
