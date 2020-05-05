package azure

import (
	"clusterCloner/poc/azure/aks_utils"
	"clusterCloner/poc/azure/aks_utils/config"
	"clusterCloner/poc/azure/resources"
	"context"
	"log"
	"os"
	"strings"
	"time"
)

var (
	aksUsername               = "azureuser"
	aksSSHPublicKeyPath       = os.Getenv("HOME") + "/.ssh/id_rsa.pub"
	aksAgentPoolCount   int32 = 1
)

func CreateClusterFromEnv(aksClusterName string) {
	var err = aks_utils.ReadEnv()
	if err != nil {
		log.Fatalf("could not set up environment: %v\n", err)
	}
	CreateCluster(config.BaseGroupName(), aksClusterName, config.DefaultLocation(), config.ClientID(), config.ClientSecret())
}

func CreateCluster(grpName string, aksClusterName string, loc string, clientID string, clientSecret string) {
	log.Printf("Group %s, Cluster %s, Location %s", grpName, aksClusterName, loc)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
	defer cancel()

	_, err := resources.CreateGroup(ctx, grpName)
	if err != nil {
		errS := err.Error()
		if strings.Contains(errS, "already exists") {
			log.Printf("Group %s already exists", grpName)
		} else {
			log.Fatal(err)
		}
	}

	//	_, err = resources.CreateGroup(ctx, "MC_joshua-playground_mycluster_westus2")
	//	if err != nil {
	//		errS := err.Error()
	//		if strings.Contains(errS, "already exists") {
	//			log.Printf("Group %s already exists", grpName)
	//		} else {
	//			log.Fatal(err)
	//		}
	//	}
	//
	_, err = createAKSCluster(ctx, aksClusterName, loc, grpName, aksUsername, aksSSHPublicKeyPath, clientID, clientSecret, aksAgentPoolCount)
	if err != nil {
		aks_utils.LogAndPanic(err)
	}
	aks_utils.PrintAndLog("created AKS cluster")
	clus, err := DescribeCluster(grpName, aksClusterName)
	_ = clus
	if err != nil {
		aks_utils.LogAndPanic(err)
	}
	aks_utils.PrintAndLog("retrieved AKS cluster")
}
