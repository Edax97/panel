package providers

import "time"

type IComServer interface {
	SendTimeValue(imei string, time time.Time, data string) (bool, error)
}

type IPanelStore interface {
	SendPanelServer(parsed [][]string, file string, serv IComServer) error
	SavePanelData(dir, file string)
}
