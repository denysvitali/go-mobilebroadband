package mb_test

import (
	"encoding/json"
	"fmt"
	mb "github.com/denysvitali/go-mobilebroadband"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStatus(t *testing.T) {
	m, err := mb.New()
	assert.Nil(t, err)
	assert.NotNil(t, m)

	st, err := m.Status()
	assert.Nil(t, err)
	assert.NotNil(t, st)
	fmt.Printf("st=%+v", st)
}

func TestGetStatus(t *testing.T) {
	m, err := mb.New()
	assert.Nil(t, err)
	assert.NotNil(t, m)

	// Pick first modem
	modems, err := m.Modems()
	assert.Nil(t, err)
	assert.Len(t, modems, 1)

	first := modems[0]
	simpleStatus, err := first.SimpleStatus()
	assert.Nil(t, err)

	jsonBytes, err := json.Marshal(simpleStatus)
	assert.Nil(t, err)

	fmt.Printf("simpleStatus: %s\n", string(jsonBytes))
}

func TestModem_GetSignal(t *testing.T) {
	m, err := mb.New()
	assert.Nil(t, err)

	modems, err := m.Modems()
	assert.Nil(t, err)
	assert.NotNil(t, modems)
	assert.Greater(t, len(modems), 0)

	modem := modems[0]
	err = modem.SetupPeriodicPolling(5)
	assert.Nil(t, err)
	time.Sleep(10)

	s, err := modem.GetSignal(mb.TechnologyLte)
	assert.Nil(t, err)
	fmt.Printf("signal=%+v", s)
	assert.NotEqual(t, "", s)
}
