package ssm

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"os"
	"context"
)

var (
	SSMParameterName = os.Getenv("SSM_PARAM")
)

func getSSMParameter(cfg aws.Config) (string, error) {
	var paramPointer = new(string)
	*paramPointer = SSMParameterName
	client := ssm.New(cfg)

	input := &ssm.GetParameterInput{
		Name: paramPointer,
	}
	// Current GetParameterRequest call can directly wrap around req, resp var. Cannot be used here
	result := client.GetParameterRequest(input)

	resp, err := result.Send(context.Background())
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ssm.ErrCodeInternalServerError:
				return "", fmt.Errorf(ssm.ErrCodeInternalServerError)
			case ssm.ErrCodeParameterNotFound:
				return "", fmt.Errorf(ssm.ErrCodeParameterNotFound)
			case ssm.ErrCodeParameterVersionNotFound:
				return "", fmt.Errorf(ssm.ErrCodeParameterVersionNotFound)
			// Can be a lot more here, did not prolong the list as these would be the most possible cases.
			}
		} else {
			return "", fmt.Errorf(err.Error())
		}

	}
	//As resp struct(obj) --> Nested struct parameter --> Value == *String
	return *resp.Parameter.Value, nil
}
