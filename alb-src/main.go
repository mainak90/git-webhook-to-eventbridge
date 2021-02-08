package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mainak90/git-webhook-to-eventbridge/cache"
	"github.com/mainak90/git-webhook-to-eventbridge/cachingclient"
	"github.com/mainak90/git-webhook-to-eventbridge/client"
	"github.com/mainak90/git-webhook-to-eventbridge/eventbus"
	"github.com/mainak90/git-webhook-to-eventbridge/validation"
	"os"
)

var (
	SecretParameterName = os.Getenv("SECRET_PARAM")
	secret = cachingclient.GetSecretCached(SecretParameterName)
)

func handle(req events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	// Parse needed values from GitHub webhook payload
	cfg := client.DefaultConfig()

	var err error

	if secret == "" {
		secret, err = cache.GenerateSecretCache(cfg, SecretParameterName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s\n", err)
			return events.ALBTargetGroupResponse{Body: err.Error(), StatusCode: 503}, nil
		}
	}

	event, delivery, signature := req.Headers["x-github-event"], req.Headers["x-github-delivery"], req.Headers["x-hub-signature"]

	if event == "" || delivery == "" {
		fmt.Fprintf(os.Stderr, "Missing x-github-event and x-hub-delivery headers")
		return events.ALBTargetGroupResponse{Body: "Missing x-github-event and x-hub-delivery headers", StatusCode: 400}, nil
	}

	if signature == "" && secret != "" {
		fmt.Fprintf(os.Stderr, "GitHub isn't providing a signature, whilst a secret is being used (please give github's webhook the secret)")
		return events.ALBTargetGroupResponse{Body: "GitHub isn't providing a signature, whilst a secret is being used (please give github's webhook the secret)", StatusCode: 400}, nil
	}

	if secret != "" {
		isValid, err := validation.IsValidPayloadSignature(secret, signature, []byte(req.Body))
		if err != nil {
			return events.ALBTargetGroupResponse{Body: err.Error(), StatusCode: 400}, nil
		}
		if (isValid) {
			fmt.Fprintf(os.Stderr, "Payload validated, coming from github...")
		}
	}

	var message interface{}
	err = json.Unmarshal([]byte(req.Body), &message)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %s\n", err)
		return events.ALBTargetGroupResponse{Body: "Unable to handle request", StatusCode: 500}, nil
	}

	fmt.Println(message)

	fmt.Println("Dispatching webhook created event..")

	err = eventbus.DispatchEvent(req, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %s\n", err)
		return events.ALBTargetGroupResponse{Body: err.Error(), StatusCode: 503}, nil
	}
	// Send response
	return events.ALBTargetGroupResponse{Body: "{ \"done\": true }", StatusCode: 200}, nil
}

func main() {
	lambda.Start(handle)
}



