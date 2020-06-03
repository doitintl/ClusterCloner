package awssdk

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/pkg/errors"
	"log"
)

// DescribeClusters ...
func DescribeClusters(region string) ([]*eks.DescribeClusterOutput, error) {

	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		printAwsErr(err)
		return nil, errors.Wrap(err, "cannot create NewSession with AWS SDK")
	}
	svc := eks.New(sess)

	input := &eks.ListClustersInput{}

	result, err := svc.ListClusters(input) //TODO use paging
	if err != nil {
		printAwsErr(err)
		return nil, errors.Wrap(err, "cannot DescribeClusters  with AWS SDK")
	}
	ret := make([]*eks.DescribeClusterOutput, 0)
	for _, clusterName := range result.Clusters {
		clusterDescription, err := DescribeCluster(*clusterName, region)

		if err != nil {
			return nil, errors.Wrap(err, "cannot DescribeCluster ("+*clusterName+")")
		}
		ret = append(ret, clusterDescription)
	}
	return ret, nil
}

// DescribeCluster ...
func DescribeCluster(clusterName, region string) (*eks.DescribeClusterOutput, error) {

	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, errors.Wrap(err, "cannot create NewSession")
	}
	svc := eks.New(sess)

	input := &eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	}

	result, err := svc.DescribeCluster(input)
	if err != nil {
		printAwsErr(err)
		return nil, errors.Wrap(err, "cannot DescribeCluster with AWS SDK")
	}

	return result, nil
}

// DescribeNodeGroups ...
func DescribeNodeGroups(clusterName string, region string) ([]*eks.DescribeNodegroupOutput, error) {

	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		printAwsErr(err)
		return nil, errors.Wrap(err, "cannot create NewSession with AWS SDK")
	}
	svc := eks.New(sess)

	input := &eks.ListNodegroupsInput{
		ClusterName: aws.String(clusterName),
	}

	result, err := svc.ListNodegroups(input) //TODO use paging
	if err != nil {
		printAwsErr(err)
		return nil, errors.Wrap(err, "cannot ListNodegroups with AWS SDK")
	}
	ret := make([]*eks.DescribeNodegroupOutput, 0)
	for _, ngName := range result.Nodegroups {
		eksNodeGroup, err := DescribeNodeGroup(clusterName, region, *ngName)
		if err != nil {
			return nil, errors.Wrap(err, "cannot DescribeCluster ("+*ngName+")")
		}
		ret = append(ret, eksNodeGroup)
	}
	return ret, nil
}

// DescribeNodeGroup ...
func DescribeNodeGroup(clusterName string, region string, ngName string) (*eks.DescribeNodegroupOutput, error) {

	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, errors.Wrap(err, "cannot create NewSession")
	}
	svc := eks.New(sess)

	input := &eks.DescribeNodegroupInput{
		ClusterName:   aws.String(clusterName),
		NodegroupName: aws.String(ngName),
	}

	result, err := svc.DescribeNodegroup(input)
	if err != nil {
		printAwsErr(err)
		return nil, errors.Wrap(err, "cannot DescribeCluster with AWS SDK")
	}

	return result, nil
}

func printAwsErr(err error) {
	if awsErr, ok := err.(awserr.Error); ok {
		log.Println(awsErr.Code(), "; ", awsErr.Error(), "; ", awsErr.Message(), "; ", awsErr.OrigErr())
	} else {
		fmt.Println(err)
	}
}
