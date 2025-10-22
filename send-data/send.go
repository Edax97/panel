package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {

	IMEI := os.Args[3]
	IP := os.Args[1]
	PORT := os.Args[2]

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	if err != nil {
		log.Panic("Conection error: ", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	login := fmt.Sprintf("2.0;%s;NA;", IMEI)
	CRC := crcChecksum([]byte(login))

	res, err := WritePacket(fmt.Sprintf("#L#%s%s\r\n", login, CRC), conn)
	if err != nil {
		log.Panic("Writing error", err)
	}
	log.Println("LOGIN:", res)

	consumption := 0
	for {
		t := time.Now()
		date := t.Format("020106")
		second := t.Format("150405")
		data := fmt.Sprintf("time:3:%s/%s,wh:1:%d;", date, second, consumption)
		message := fmt.Sprintf("%s;%s.000000000;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;;NA;%s", date, second, data)
		//hourStr := time.Now().Format("2006.01.02.15.04")
		//consumption += 1000
		//data := fmt.Sprintf("time:3:%s;WHrecibidos:1:%d;", hourStr, consumption)
		//message := fmt.Sprintf("NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;;NA;%s", data)
		CRC = crcChecksum([]byte(message))
		res, err = WritePacket(fmt.Sprintf("#D#%s%s\r\n", message, CRC), conn)
		if err != nil {
			log.Println("Writing error", err)
		}
		log.Println("DATA:", res)
		time.Sleep(time.Second * 10)
	}
}
