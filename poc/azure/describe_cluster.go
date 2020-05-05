package azure

import (
	"clusterCloner/poc/azure/aks_utils"
	"clusterCloner/poc/utilities"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2017-09-30/containerservice"
	"log"
	"time"
)

// describeCluster returns an existing AKS cluster given a resource group name and resource name
func DescribeCluster(grpName, clusterName string) (c containerservice.ManagedCluster, err error) {
	err_ := aks_utils.ReadEnv()
	_ = err_
	log.Printf("Group %s, Cluster %s", grpName, clusterName)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
	defer cancel()
	aksClient, err := getAKSClient()
	if err != nil {
		return c, fmt.Errorf("cannot get AKS client: %v", err)
	}
	c, err = aksClient.Get(ctx, grpName, clusterName)
	var list containerservice.ManagedClusterListResultPage

	log.Print(list)
	if err != nil {
		return c, fmt.Errorf("cannot get AKS managed cluster %v from resource group %v: %v", clusterName, grpName, err)
	}
	//	props := c.ManagedClusterProperties
	utilities.PrintAsJson(c)
	return c, nil
}
