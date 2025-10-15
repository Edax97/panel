package wailonServer

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type WailonServer struct {
	ip   string
	port string
}

func NewWailonServer(ip string, port string) *WailonServer {
	fmt.Println("Server started:", ip, port)
	return &WailonServer{ip, port}
}

func (s *WailonServer) SendTimeValue(imei string, date time.Time, value int) (bool, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", s.ip, s.port))
	if err != nil {
		return false, err
	}
	defer func() {
		_ = conn.Close()
	}()

	login := fmt.Sprintf("2.0;%s;NA;", imei)
	CRC := crcChecksum([]byte(login))

	_, err = writePacket(fmt.Sprintf("#L#%s%s\r\n", login, CRC), conn)
	if err != nil {
		return false, err
	}
	//fmt.Println("- Login: ", res)

	//ddmmyy;hhmmss
	datetimeStr := date.Format("2006-01-02/15-04-05")
	data := fmt.Sprintf("time:3:%s,wh:1:%d;", datetimeStr, value)
	message := fmt.Sprintf("NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;;NA;%s", data)
	CRC = crcChecksum([]byte(message))
	//fmt.Println("- Data: ", message)
	res, err := writePacket(fmt.Sprintf("#D#%s%s\r\n", message, CRC), conn)
	if err != nil {
		return false, err
	}
	//#AD#1
	status := strings.Split(res+"##", "#")[2]
	if status != "1" {
		fmt.Println("Error sending data:", status)
		return false, fmt.Errorf("IMEI %s, (%s, %s)", imei, message, res)
	}
	//fmt.Println("- Res: ", res)
	return true, nil
}
