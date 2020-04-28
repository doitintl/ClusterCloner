package eks

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
)

func CreateCluster(name string) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-west-2")})
	svc := eks.New(sess)

	input := &eks.CreateClusterInput{
		ClientRequestToken: aws.String("4d2120a1-3d38-460a-9756-e6b97fddb955"),
		Name:               aws.String(name),
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			//			SecurityGroupIds: []*string{
			//				aws.String("sg-6979fe18"),
			//			},
			SubnetIds: []*string{
				aws.String("subnet-53c2272b"),
				aws.String("subnet-5800f405"),
			},
		},
		RoleArn: aws.String("arn:aws:iam::649592902942:role/eksClusterRole"),
		Version: aws.String("1.15"),
	}

	result, err := svc.CreateCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceInUseException:
				fmt.Println(eks.ErrCodeResourceInUseException, aerr.Error())
			case eks.ErrCodeResourceLimitExceededException:
				fmt.Println(eks.ErrCodeResourceLimitExceededException, aerr.Error())
			case eks.ErrCodeInvalidParameterException:
				fmt.Println(eks.ErrCodeInvalidParameterException, aerr.Error())
			case eks.ErrCodeClientException:
				fmt.Println(eks.ErrCodeClientException, aerr.Error())
			case eks.ErrCodeServerException:
				fmt.Println(eks.ErrCodeServerException, aerr.Error())
			case eks.ErrCodeServiceUnavailableException:
				fmt.Println(eks.ErrCodeServiceUnavailableException, aerr.Error())
			case eks.ErrCodeUnsupportedAvailabilityZoneException:
				fmt.Println(eks.ErrCodeUnsupportedAvailabilityZoneException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}
