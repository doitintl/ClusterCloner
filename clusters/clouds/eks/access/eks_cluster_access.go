package access

import (
	"clusterCloner/clusters/cluster_info"
	"clusterCloner/clusters/util"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"log"
)

type EksClusterAccess struct {
}

func (ca EksClusterAccess) CreateCluster(info cluster_info.ClusterInfo) (cluster_info.ClusterInfo, error) {
	panic("implement me")
}
func (ca EksClusterAccess) DescribeCluster(clusterName string, region string) (cluster_info.ClusterInfo, error) {

	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	svc := eks.New(sess)

	input := &eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	}

	result, err := svc.DescribeCluster(input)
	if err != nil {
		printAwsErr(err)
		return cluster_info.ClusterInfo{}, err
	}
	log.Println(util.MarshallToJsonString(result))
	return cluster_info.ClusterInfo{Scope: "", Location: region, Name: clusterName, NodeCount: 1, GeneratedBy: cluster_info.READ}, nil
}

func (ca EksClusterAccess) ListClusters(_ string, location string) ([]cluster_info.ClusterInfo, error) {
	clusterNames, err := clusterNames(location)
	if err != nil {
		return nil, err
	}
	ret := make([]cluster_info.ClusterInfo, 0)

	for _, s := range clusterNames {
		log.Print(*s)
		var clusterInfo, err_ = ca.DescribeCluster(location, *s)
		if err_ != nil {
			log.Print("Error ", err_)
			return nil, err_
		}
		ret = append(ret, clusterInfo)
	}
	return ret, nil
}

func clusterNames(region string) ([]*string, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	input := &eks.ListClustersInput{}
	svc := eks.New(sess)

	result, err := svc.ListClusters(input)
	if err != nil {
		printAwsErr(err)
		return nil, err
	}
	var clusterNames = result.Clusters
	return clusterNames, nil
}

func printAwsErr(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case eks.ErrCodeInvalidParameterException:
			log.Println(eks.ErrCodeInvalidParameterException, aerr.Error())
		case eks.ErrCodeClientException:
			log.Println(eks.ErrCodeClientException, aerr.Error())
		case eks.ErrCodeServerException:
			log.Println(eks.ErrCodeServerException, aerr.Error())
		case eks.ErrCodeServiceUnavailableException:
			log.Println(eks.ErrCodeServiceUnavailableException, aerr.Error())
		default:
			log.Println(aerr.Error())
		}
	} else {
		// cast err to awserror.Error to get the Code and Message from an error.
		awserror := err.(awserr.Error)
		fmt.Println(awserror.Message())
	}
}
