package wailonServer

import (
	"bufio"
	"net"
)

func writePacket(packet string, con net.Conn) (string, error) {
	_, err := con.Write([]byte(packet))
	if err != nil {
		return "", err
	}
	reader := bufio.NewReader(con)
	res, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return res, nil
}
