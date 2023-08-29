package config

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	cfg aws.Config
)

func init() {
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	cfg = awsConfig
}

func GetAWSConfig() aws.Config {
	return cfg
}
