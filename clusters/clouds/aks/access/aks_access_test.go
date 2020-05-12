package access

import (
	"context"
	"log"
	"testing"
	"time"
)

func init() {
	supportedVersions = []string{"1.14.9", "1.14.8", "1.14.11", "1.15.1", "1.15.8"}
}

func TestDescribeCluster(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
	defer cancel()
	clus, _ := getCluster(ctx, "joshua-playground", "cluster-2-1havh-paiuq")
	log.Println(clus)
}
func TestParseMachineType(t *testing.T) {
	machineType := "Standard_D1_v2"

	mt := MachineTypeByName(machineType)
	if mt.Name != machineType {
		t.Error(mt.Name)
	}
	if mt.CPU != 1 {
		t.Error(mt.CPU)

	}
	if mt.RAMMB != 3584 {
		t.Error(mt.RAMMB)
	}
}
