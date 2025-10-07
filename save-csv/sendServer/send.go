package sendServer

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var IP string = os.Getenv("SERVER_IP")
var PORT string = os.Getenv("SERVER_PORT")

func SendMessage(imei string, date time.Time, consumption int) (bool, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	if err != nil {
		return false, err
	}
	defer func() {
		_ = conn.Close()
	}()

	login := fmt.Sprintf("2.0;%s;NA;", imei)
	CRC := crcChecksum([]byte(login))

	res, err := writePacket(fmt.Sprintf("#L#%s%s\r\n", login, CRC), conn)
	if err != nil {
		return false, err
	}
	log.Println("LOGIN:", res)

	hourStr := date.Format("2006.01.02.15.04")
	data := fmt.Sprintf("time:3:%s;WHrecibidos:1:%d;", hourStr, consumption)
	message := fmt.Sprintf("NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;;NA;%s", data)
	CRC = crcChecksum([]byte(message))
	res, err = writePacket(fmt.Sprintf("#D#%s%s\r\n", message, CRC), conn)
	if err != nil {
		return false, err
	}
	log.Println("DATA:", res)
	return true, nil
}
