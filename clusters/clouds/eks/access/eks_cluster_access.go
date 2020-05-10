package access

import (
	"clustercloner/clusters/clusterinfo"
	"clustercloner/clusters/util"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"log"
)

// EKSClusterAccess ...
type EKSClusterAccess struct {
}

// CreateCluster ...
func (ca EKSClusterAccess) CreateCluster(info clusterinfo.ClusterInfo) (clusterinfo.ClusterInfo, error) {
	panic("implement me")
}

// DescribeCluster ...
func (ca EKSClusterAccess) DescribeCluster(clusterName string, region string) (clusterinfo.ClusterInfo, error) {

	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	svc := eks.New(sess)

	input := &eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	}

	result, err := svc.DescribeCluster(input)
	if err != nil {
		printAwsErr(err)
		return clusterinfo.ClusterInfo{}, err
	}
	log.Println(util.MarshallToJSONString(result))
	return clusterinfo.ClusterInfo{Scope: "", Location: region, Name: clusterName, GeneratedBy: clusterinfo.READ}, nil
}

// ListClusters ...
func (ca EKSClusterAccess) ListClusters(_ string, location string) ([]clusterinfo.ClusterInfo, error) {
	clusterNames, err := clusterNames(location)
	if err != nil {
		return nil, err
	}
	ret := make([]clusterinfo.ClusterInfo, 0)

	for _, s := range clusterNames {
		log.Print(*s)
		var clusterInfo, err = ca.DescribeCluster(location, *s)
		if err != nil {
			log.Print("Error ", err)
			return nil, err
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
