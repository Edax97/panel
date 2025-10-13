package wailonServer

import (
	"fmt"
	"time"
)

type MockServer struct {
}

func NewMockServer() *MockServer {
	return &MockServer{}
}

func (s *MockServer) SendTimeValue(imei string, date time.Time, value int) (bool, error) {
	login := fmt.Sprintf("2.0;%s;NA;", imei)
	CRC := crcChecksum([]byte(login))
	loginPacket := fmt.Sprintf("#L#%s%s\r\n", login, CRC)
	fmt.Printf("  - IMEI: %s\n    LOGIN PACKET: %s\n", imei, loginPacket)

	hourStr := date.Format("2006.01.02.15.04")
	data := fmt.Sprintf("time:3:%s;WHrecibidos:1:%d;", hourStr, value)
	message := fmt.Sprintf("NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;;NA;%s", data)
	CRC = crcChecksum([]byte(message))
	dataPacket := fmt.Sprintf("#D#%s%s\r\n", message, CRC)
	fmt.Printf("  - DATA PACKET: %s\n", dataPacket)

	return true, nil
}
