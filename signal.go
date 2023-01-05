package mb

import (
	"fmt"
	"github.com/godbus/dbus/v5"
)

type Signal struct {
	Rsrp float64 `json:"rsrp"`
	Rsrq float64 `json:"rsrq"`
	Rssi float64 `json:"rssi"`
	Snr float64 `json:"snr"`
}

func (m Modem) GetSignal(technology Technology) (*Signal, error) {
	var dbusResult map[string]dbus.Variant
	p, err := m.obj.GetProperty(fmt.Sprintf("org.freedesktop.ModemManager1.Modem.Signal.%s", technology))
	if err != nil {
		return nil, err
	}
	err = p.Store(&dbusResult)

	if len(dbusResult) == 0 {
		return nil, nil
	}

	s := Signal{}
	if v, ok := dbusResult["rsrp"]; ok {
		s.Rsrp = v.Value().(float64)
	}

	if v, ok := dbusResult["rsrq"]; ok {
		s.Rsrq = v.Value().(float64)
	}

	if v, ok := dbusResult["rssi"]; ok {
		s.Rssi = v.Value().(float64)
	}

	if v, ok := dbusResult["snr"]; ok {
		s.Snr = v.Value().(float64)
	}

	return &s, err
}