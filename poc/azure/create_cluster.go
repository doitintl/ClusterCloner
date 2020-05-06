package azure

import (
	"clustercloner/poc/azure/aksutils"
	"clustercloner/poc/azure/aksutils/config"
	"clustercloner/poc/azure/resources"
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

// CreateClusterFromEnv ...
func CreateClusterFromEnv(aksClusterName string) {
	var err = aksutils.ReadEnv()
	if err != nil {
		log.Fatalf("could not set up environment: %v\n", err)
	}
	CreateCluster(config.BaseGroupName(), aksClusterName, config.DefaultLocation(), config.ClientID(), config.ClientSecret())
}

// CreateCluster ...
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
		aksutils.LogAndPanic(err)
	}
	aksutils.PrintAndLog("created AKS cluster")
	clus, err := DescribeCluster(grpName, aksClusterName)
	_ = clus
	if err != nil {
		aksutils.LogAndPanic(err)
	}
	aksutils.PrintAndLog("retrieved AKS cluster")
}
