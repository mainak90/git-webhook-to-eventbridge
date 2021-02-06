package cache_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/secretsmanageriface"
	"github.com/aws/aws-sdk-go/aws"
	"testing"
)

type MockSecretsManagerClient struct {
	secretsmanageriface.ClientAPI
}

func (m *MockSecretsManagerClient) getSecretValueRequest(input *secretsmanager.GetSecretValueInput) (resp secretsmanager.GetSecretValueOutput) {

	var res secretsmanager.GetSecretValueOutput
	var dummyresp = new(string)
	*dummyresp = "this is the dummy value"
	var name = new(string)
	*name = "dummy_secret"
	res.SecretString = dummyresp
	res.Name = name
	return res
}

func TestGenerateSecretCache(t *testing.T) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String("dummy_secret"),
		VersionStage: aws.String("AWSCURRENT"),
	}
	var m MockSecretsManagerClient
	res := m.GetSecretValueRequest(input)
	resp, err := res.Send(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	if *resp.SecretString != "this is the dummy value" {
		t.Errorf("Unexpected value %s", *resp.SecretString)
	}
}




