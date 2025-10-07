package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var imeiMap = make(map[string]string)
var imeiList = make([][]string, 1)

func main() {
	var deviceStr string
	imeiList[0] = []string{"ID", "IMEI"}
	InputFile := os.Args[1]
	ImeiFile := os.Args[2]

	reg := regexp.MustCompile("[^a-zA-Z0-9]+") // Matches any character not a letter or number

	fileData, err := readCSV(InputFile, ',')
	if err != nil {
		log.Panic(err)
	}
	deviceIds := fileData[1]
	for _, device := range deviceIds {
		_, ok := imeiMap[device]
		if !ok {
			//base convert decoding
			deviceSplit := strings.Split(device, ":")
			if len(deviceSplit) < 2 {
				continue
			}
			deviceStr = reg.ReplaceAllString(deviceSplit[1], "")
			imei, err := strconv.ParseInt(deviceStr, 36, 64)
			if err != nil {
				log.Panic(err)
			}
			imei = 100000000000000 + imei%100000000000000
			imeiStr := strconv.FormatInt(imei, 10)
			imeiMap[device] = imeiStr
			imeiList = append(imeiList, []string{device, imeiStr})
		}
	}
	saveSCV(ImeiFile, imeiList)
}
