package wailonServer

import (
	"fmt"
	"log"
	"net"
	"time"
)

type WailonServer struct {
	ip   string
	port string
}

func NewWailonServer(ip string, port string) *WailonServer {
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

	res, err := writePacket(fmt.Sprintf("#L#%s%s\r\n", login, CRC), conn)
	if err != nil {
		return false, err
	}
	log.Println("LOGIN:", res)

	hourStr := date.Format("2006.01.02.15.04")
	data := fmt.Sprintf("time:3:%s;WHrecibidos:1:%d;", hourStr, value)
	message := fmt.Sprintf("NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;;NA;%s", data)
	CRC = crcChecksum([]byte(message))
	res, err = writePacket(fmt.Sprintf("#D#%s%s\r\n", message, CRC), conn)
	if err != nil {
		return false, err
	}
	log.Println("DATA:", res)
	return true, nil
}
