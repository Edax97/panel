package wailonServer

import (
	"fmt"
	"net"
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

	res, err := writePacket(fmt.Sprintf("#L#%s%s\r\n", login, CRC), conn)
	if err != nil {
		return false, err
	}
	fmt.Println("- Login: ", res)

	//ddmmyy;hhmmss
	dateStr := date.UTC().Format("020106")
	timeStr := date.UTC().Format("150405") + ".000000000"
	datetimeStr := date.Format("2006.01.02.15.04.05")
	data := fmt.Sprintf("time:3:%s,wh:1:%d;", datetimeStr, value)
	message := fmt.Sprintf("%s;%s;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;;NA;%s", dateStr, timeStr, data)
	CRC = crcChecksum([]byte(message))
	fmt.Println("- Data: ", message)
	res, err = writePacket(fmt.Sprintf("#D#%s%s\r\n", message, CRC), conn)
	if err != nil {
		return false, err
	}
	fmt.Println("- Res: ", res)
	return true, nil
}
