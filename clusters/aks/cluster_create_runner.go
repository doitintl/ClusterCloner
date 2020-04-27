package aks

import (
	"clusterCloner/clusters/aks/resources"
	"clusterCloner/clusters/aks/utils/config"
	"clusterCloner/clusters/aks/utils/util"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

var (
	sshPublicKeyPath = os.Getenv("HOME") + "/.ssh/id_rsa.pub"

	containerGroupName        = "gosdk-aci"
	aksClusterName            = "gosdk-aks"
	aksUsername               = "azureuser"
	aksSSHPublicKeyPath       = os.Getenv("HOME") + "/.ssh/id_rsa.pub"
	aksAgentPoolCount   int32 = 1
)

func addEnv() error {
	err := config.ParseEnvironment()
	if err != nil {
		return fmt.Errorf("failed to add top-level env: %+v", err)
	}
	return nil
}

func setup() error {
	var err error
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	err = addEnv()
	if err != nil {
		return err
	}

	return nil
}

func CreateCluster() {

	var err = setup()
	if err != nil {
		log.Fatalf("could not set up environment: %v\n", err)
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
	defer cancel()

	_, err = resources.CreateGroup(ctx, config.BaseGroupName())
	if err != nil {
		log.Printf("Group %s already exists\n", config.BaseGroupName())
	}

	loc := config.DefaultLocation()
	_, err = createAKS(ctx, aksClusterName, loc, config.BaseGroupName(), aksUsername, aksSSHPublicKeyPath, config.ClientID(), config.ClientSecret(), aksAgentPoolCount)
	if err != nil {
		util.LogAndPanic(err)
	}
	util.PrintAndLog("created AKS cluster")

	_, err = getAKS(ctx, config.BaseGroupName(), aksClusterName)
	if err != nil {
		util.LogAndPanic(err)
	}
	util.PrintAndLog("retrieved AKS cluster")
}
