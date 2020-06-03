package access

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMachineTypes(t *testing.T) {
	machineTypeCount := GetMachineTypes().Length()
	assert.Greater(t, machineTypeCount, 250)
	assert.Less(t, machineTypeCount, 350) //check the filter for only EKS-supported machine types
	mt, err := GetMachineTypes().Get("m5.12xlarge")
	assert.Nil(t, err)
	assert.Equal(t, "m5.12xlarge", mt.Name)
	assert.Equal(t, 196608, mt.RAMMB)
	assert.Equal(t, 48, mt.CPU)
}
