package machinetypes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMachineTypes(t *testing.T) {
	var m = NewMachineTypeMap()
	m.Set("a", MachineType{
		Name:  "a",
		CPU:   10,
		RAMMB: 11,
	})
	m.Set("b", MachineType{
		Name:  "b",
		CPU:   20,
		RAMMB: 21,
	})

	m.Set("a", MachineType{
		Name:  "a",
		CPU:   30,
		RAMMB: 31,
	})
	mta, err := m.Get("b")
	assert.Nil(t, err)
	assert.Equal(t, 20, mta.CPU)
	mta, err = m.Get("a")
	assert.Nil(t, err)
	assert.Equal(t, 30, mta.CPU)
	mta, err = m.Get("not-present")
	assert.NotNil(t, err)
	assert.Equal(t, "", mta.Name)

}
