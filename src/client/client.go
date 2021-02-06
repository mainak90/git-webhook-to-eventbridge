package client

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go/aws/endpoints"
)

func DefaultConfig() aws.Config {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load sdk config, " + err.Error())
	}
	cfg.Region = endpoints.UsEast1RegionID
	return cfg
}
