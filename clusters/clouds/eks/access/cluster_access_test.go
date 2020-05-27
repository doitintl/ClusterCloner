package access

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMachineTypes(t *testing.T) {
	machineTypeCount := len(MachineTypes)
	assert.Greater(t, machineTypeCount, 20)
	mt := MachineTypeByName("m4.4xlarge")
	assert.Equal(t, mt.Name, "m4.4xlarge")
	assert.Equal(t, mt.RAMMB, 65536)
	assert.Equal(t, mt.CPU, 16)
}
