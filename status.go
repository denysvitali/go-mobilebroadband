package mb

import (
	"encoding/json"
	"fmt"
	"github.com/godbus/dbus/v5"
	"strings"
)

const ModemManager1 = "/org/freedesktop/ModemManager1"
const ModemManager1Dest = "org.freedesktop.ModemManager1"

type MobileBroadband struct {
	dbus *dbus.Conn
}

type Status struct {
	Modem dbus.BusObject
	M3gpp Modem3gpp
}

var _ json.Unmarshaler = (*SignalQuality)(nil)

type SignalQuality struct {
	Value             float64
	RecentlyRefreshed bool
}

func (s *SignalQuality) UnmarshalJSON(bytes []byte) error {
	var v []interface{}
	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	s.Value = v[0].(float64)
	s.RecentlyRefreshed = v[1].(bool)
	return nil
}

type Cdma1xRegistrationState = int
type CdmaEvdoRegistrationState = int

type SimpleStatus struct {
	State                     uint32                    `json:"state"`
	SignalQuality             SignalQuality             `json:"signal-quality"`
	CurrentBands              []uint32                  `json:"current-bands"`
	AccessTechnologies        int                       `json:"access-technologies"`
	RegistrationState         RegistrationState         `json:"registration-state"`
	OperatorCode              string                    `json:"m3gpp-operator-code"`
	OperatorName              string                    `json:"m3gpp-operator-name"`
	Cdma1xRegistrationState   Cdma1xRegistrationState   `json:"cdma-cdma1x-registration-state"`
	CdmaEvdoRegistrationState CdmaEvdoRegistrationState `json:"cdma-evdo-registration-state"`
	CdmaSid                   uint32                    `json:"cdma-sid"`
	CdmaNid                   uint32                    `json:"cdma-nid"`
}

type RegistrationState int

type Modem3gpp struct {
	Imei              string
	RegistrationState RegistrationState
	OperatorCode      string
	OperatorName      string
}

func New() (*MobileBroadband, error) {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, fmt.Errorf("unable to create DBUS Session: %v", err)
	}

	return &MobileBroadband{
		dbus: conn,
	}, nil
}

func (m *MobileBroadband) Modems() ([]Modem, error) {
	var objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	obj := m.dbus.Object(ModemManager1Dest, ModemManager1)
	err := obj.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&objects)
	if err != nil {
		return nil, err
	}

	var modems []Modem
	for k, _ := range objects {
		if strings.HasPrefix(string(k), fmt.Sprintf("%s/Modem/", ModemManager1)) {
			modems = append(modems, Modem{obj: m.dbus.Object(ModemManager1Dest, k)})
		}
	}
	return modems, nil
}

func (m *MobileBroadband) Status() ([]Status, error) {
	// Get Modems
	modems, err := m.Modems()
	if err != nil {
		return nil, err
	}

	var statuses []Status
	for _, modem := range modems {
		status := Status{}
		status.Modem = modem.obj
		status.M3gpp = Modem3gpp{
			Imei:              modem.getImei(),
			RegistrationState: modem.getRegistrationState(),
			OperatorCode:      modem.getOperatorCode(),
			OperatorName:      modem.getOperatorName(),
		}
		signals := []string{"Lte", "Nr5g", "Umts", "Cdma", "Rate"}
		for _, s := range signals {
			fmt.Printf("%s: %s\n", s, modem.getSignal(s))
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

type Modem struct {
	obj dbus.BusObject
}

func (m Modem) getImei() string {
	p, err := m.obj.GetProperty("org.freedesktop.ModemManager1.Modem.Modem3gpp.Imei")
	if err != nil {
		return ""
	}
	switch t := p.Value().(type) {
	case string:
		return t
	}
	return ""
}

func (m Modem) getRegistrationState() RegistrationState {
	p, err := m.obj.GetProperty("org.freedesktop.ModemManager1.Modem.Modem3gpp.RegistrationState")
	if err != nil {
		return 0
	}
	switch t := p.Value().(type) {
	case uint32:
		return RegistrationState(t)
	}
	return 0
}

func (m Modem) getOperatorCode() string {
	p, err := m.obj.GetProperty("org.freedesktop.ModemManager1.Modem.Modem3gpp.OperatorCode")
	if err != nil {
		return ""
	}
	switch t := p.Value().(type) {
	case string:
		return t
	}
	return ""
}

func (m Modem) SimpleStatus() (*SimpleStatus, error) {
	var status map[string]interface{}
	var statusOutput SimpleStatus
	err := m.obj.Call("org.freedesktop.ModemManager1.Modem.Simple.GetStatus", 0).Store(&status)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(status)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &statusOutput)
	if err != nil {
		return nil, err
	}

	return &statusOutput, nil
}

func (m Modem) getSignal(technology string) string {
	var signal interface{}
	p, err := m.obj.GetProperty(fmt.Sprintf("org.freedesktop.ModemManager1.Modem.Signal.%s", technology))
	if err != nil {
		return ""
	}
	signal = p.Value()

	return fmt.Sprintf("%+v", signal)
}

func (m Modem) getOperatorName() string {
	p, err := m.obj.GetProperty("org.freedesktop.ModemManager1.Modem.Modem3gpp.OperatorName")
	if err != nil {
		return ""
	}
	switch t := p.Value().(type) {
	case string:
		return t
	}
	return ""
}
