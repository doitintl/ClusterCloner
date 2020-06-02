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
			return nil, errors.Wrap(err, "cannot DescribeNodeGroup ("+*ngName+")")
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
		return nil, errors.Wrap(err, "cannot DescribeNodeGroup with AWS SDK")
	}

	return result, nil
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
