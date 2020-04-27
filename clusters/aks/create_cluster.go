package aks

import (
	"clusterCloner/clusters/aks/utils"
	"clusterCloner/clusters/aks/utils/config"
	"log"
	"os"
)

var (
	aksUsername               = "azureuser"
	aksSSHPublicKeyPath       = os.Getenv("HOME") + "/.ssh/id_rsa.pub"
	aksAgentPoolCount   int32 = 1
)

func CreateClusterFromEnv(aksClusterName string) {
	var err = utils.ReadEnv()
	if err != nil {
		log.Fatalf("could not set up environment: %v\n", err)
	}
	CreateCluster(config.BaseGroupName(), aksClusterName, config.DefaultLocation(), config.ClientID(), config.ClientSecret())
}

func CreateCluster(grpName string, aksClusterName string, loc string, clientID string, clientSecret string) {
	//	log.Printf("Group %s, Cluster %s, Location %s", grpName, aksClusterName, loc)
	//	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
	//	defer cancel()
	//
	//	_, err := resources.CreateGroup(ctx, grpName)
	//	if err != nil {
	//		errS := err.Error()
	//		if strings.Contains(errS, "already exists") {
	//			log.Printf("Group %s already exists", grpName)
	//		} else {
	//			log.Fatal(err)
	//		}
	//	}
	//
	//	_, err = createAKSCluster(ctx, aksClusterName, loc, grpName, aksUsername, aksSSHPublicKeyPath, clientID, clientSecret, aksAgentPoolCount)
	//	if err != nil {
	//		utils.LogAndPanic(err)
	//	}
	//	utils.PrintAndLog("created AKS cluster")
	clus, err := ReadCluster(grpName, aksClusterName)
	_ = clus
	if err != nil {
		utils.LogAndPanic(err)
	}
	utils.PrintAndLog("retrieved AKS cluster")
}
