package wailonServer

import (
	"fmt"
	"sync"
	"time"
)

type MockServer struct {
	OutBuffer string
	mutex     sync.Mutex
}

func NewMockServer() *MockServer {
	return &MockServer{}
}

func (s *MockServer) SendTimeValue(imei string, date time.Time, data string) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	login := fmt.Sprintf("2.0;%s;NA;", imei)
	CRC := crcChecksum([]byte(login))
	loginPacket := fmt.Sprintf("#L#%s%s\r\n", login, CRC)
	s.OutBuffer = fmt.Sprintf("%s\n  - IMEI: %s\n    LOGIN PACKET: %s\n", s.OutBuffer, imei, loginPacket)

	hourStr := date.In(time.UTC).Format("2006.01.02.15.04")
	message := fmt.Sprintf("%s;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;;NA;%s", hourStr, data)
	CRC = crcChecksum([]byte(message))
	dataPacket := fmt.Sprintf("#D#%s%s\r\n", message, CRC)
	s.OutBuffer = fmt.Sprintf("%s\n - DATA PACKET: %s\n", s.OutBuffer, dataPacket)

	fmt.Println(s.OutBuffer)
	return true, nil
}

func (s *MockServer) FlushOut() {
	s.OutBuffer = ""
}
