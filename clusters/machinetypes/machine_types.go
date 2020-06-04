package machinetypes

import (
	"fmt"
	"github.com/pkg/errors"
)

// MachineType ...
type MachineType struct {
	Name  string
	CPU   int
	RAMMB int
}

// NewMachineTypeMap ...
func NewMachineTypeMap() MachineTypes {
	ret := MachineTypes{}
	ret.mts = make([]MachineType, 0)
	return ret
}

// MachineTypes ...
type MachineTypes struct {
	mts []MachineType
}

// Get ...
func (m *MachineTypes) Get(key string) (machineType MachineType, err error) {
	for _, mt := range m.mts {
		if mt.Name == key {
			if machineType.Name != "" {
				return MachineType{}, errors.New("multiple occurences of " + key)
			}
			machineType = mt
		}
	}
	if machineType.Name == "" {
		return MachineType{}, errors.New(key + " not found")
	}
	return machineType, nil
}

// Set ...
func (m *MachineTypes) Set(key string, value MachineType) {
	found := -1
	for idx, mt := range m.mts {
		if mt.Name == key {
			if found != -1 {
				panic(fmt.Sprintf("%s appeared at both %d and %d", key, idx, found))
			}
			found = idx
		}
	}
	if found != -1 {
		m.mts[found] = value //replace value
	} else {
		m.mts = append(m.mts, value)
	}

}

// Length ...
func (m *MachineTypes) Length() int {
	return len(m.mts)

}

// List ...
func (m *MachineTypes) List() []MachineType {
	return m.mts
}
