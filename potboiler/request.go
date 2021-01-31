package potboiler

import (
	"errors"
	"github.com/mainak90/git-webhook-to-eventbridge"

	"github.com/aws/aws-lambda-go/events"
)

// Request contains data from API Gateway invocation parameters
type Request struct {
	Channel string
	Payload main.Payload
}

func parseRequest(req events.APIGatewayProxyRequest) (*Request, error) {
	payload, err := main.parsePayload([]byte(req.Body))

	if err != nil {
		return nil, errors.New("Unable to parse request payload")
	}

	return &Request{
		req.PathParameters["channel"],
		*payload,
	}, nil
}
