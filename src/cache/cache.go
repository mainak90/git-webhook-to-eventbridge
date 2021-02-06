package cache

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GenerateSecretCache(cfg aws.Config, secretname string) (string, error) {
	client := secretsmanager.New(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretname),
		VersionStage: aws.String("AWSCURRENT"),
	}
	result := client.GetSecretValueRequest(input)
	resp, err := result.Send(context.Background())
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
				case secretsmanager.ErrCodeDecryptionFailure:
					return "", fmt.Errorf(secretsmanager.ErrCodeDecryptionFailure)
				case secretsmanager.ErrCodeInternalServiceError:
					return "", fmt.Errorf(secretsmanager.ErrCodeInternalServiceError)
				case secretsmanager.ErrCodeInvalidParameterException:
					return "", fmt.Errorf(secretsmanager.ErrCodeInvalidParameterException)
				case secretsmanager.ErrCodeInvalidRequestException:
					return "", fmt.Errorf(secretsmanager.ErrCodeInvalidRequestException)
				case secretsmanager.ErrCodeResourceNotFoundException:
					return "", fmt.Errorf(secretsmanager.ErrCodeResourceNotFoundException)
			}
		} else {
			return "", fmt.Errorf(err.Error())
		}

	}
	return *resp.SecretString, nil
}