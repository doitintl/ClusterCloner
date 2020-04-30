package aks

import (
	"clusterCloner/clusters"
	"clusterCloner/clusters/aks/config"
	"clusterCloner/clusters/aks/iam"
	"clusterCloner/clusters/util"
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice"
	_ "log"
	_ "reflect"
	"time"
)

func init() {
	_ = util.ReadEnv()
}

type AksClusterAccess struct {
}

func (ca AksClusterAccess) ListClusters(subscription string, location string) (ci []clusters.ClusterInfo, err error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
	defer cancel()
	var aksClient, err_ = getAKSClient()
	if err_ != nil {
		return ci, errors.New("cannot get AKS client")
	}
	ret := make([]clusters.ClusterInfo, 0)

	clusterList, _ := aksClient.List(ctx)
	for _, managedCluster := range clusterList.Values() {
		var props = managedCluster.ManagedClusterProperties

		var count int32 = 0
		for _, app := range *props.AgentPoolProfiles {
			count += *app.Count
		}
		ci := clusters.ClusterInfo{Scope: subscription, Location: location, Name: *managedCluster.Name, NodeCount: count}
		ret = append(ret, ci)

	}
	return ret, nil
}

func getAKSClient() (mcc containerservice.ManagedClustersClient, err error) {
	aksClient := containerservice.NewManagedClustersClient(config.SubscriptionID())
	auth, err := iam.GetResourceManagementAuthorizer()
	if err != nil {
		return mcc, err
	}
	aksClient.Authorizer = auth
	_ = aksClient.AddToUserAgent(config.UserAgent())
	aksClient.PollingDuration = time.Hour * 1
	return aksClient, nil
}