package mb_test

import (
	"fmt"
	mb "github.com/denysvitali/go-mobilebroadband"
	"github.com/stretchr/testify/assert"
	"testing"
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

	fmt.Printf("simpleStatus: %+v\n", simpleStatus)
}
