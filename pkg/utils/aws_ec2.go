package utils

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func (u *Utils) GetInstancesDetails(cfg aws.Config, params *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	client := ec2.NewFromConfig(cfg)
	return client.DescribeInstances(context.TODO(), params)
}
