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

func (s *WailonServer) SendTimeValue(imei string, t time.Time, value int) (bool, error) {
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
	fmt.Println("- Login: ", login)

	//ddmmyy;hhmmss
	date := t.Format("020106")
	second := t.Format("150405")
	data := fmt.Sprintf("time:3:%s/%s,wh:1:%d;", date, second, value)
	message := fmt.Sprintf("%s;%s.000000000;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;NA;;NA;%s", date, second, data)
	CRC = crcChecksum([]byte(message))
	fmt.Println("- Data: ", message)
	res, err = writePacket(fmt.Sprintf("#D#%s%s\r\n", message, CRC), conn)
	fmt.Println("- Res: ", res)
	if err != nil {
		return false, err
	}
	//#AD#1
	if !strings.Contains(res, "#AD#1") {
		return false, fmt.Errorf("IMEI %s, (%s, %s)", imei, message, res)
	}
	return true, nil
}
