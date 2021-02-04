package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/mainak90/git-webhook-to-eventbridge/validation"
	"os"
)

var (
	EventBusName = os.Getenv("EVENT_BUS_NAME")
	EventSource = "gitwebhook.lambda"
	EventDetail = os.Getenv("EVENT_DETAIL")
	SSMParameterName = os.Getenv("SSM_PARAM")
)

func defaultConfig() aws.Config {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load sdk config, " + err.Error())
	}
	cfg.Region = endpoints.EuWest1RegionID
	return cfg
}

func eventBridgeSession(cfg aws.Config) *eventbridge.Client{
	return eventbridge.New(cfg)
}

func eventRequestEntry(details string) eventbridge.PutEventsInput{
	return eventbridge.PutEventsInput{Entries: []eventbridge.PutEventsRequestEntry{
		{
			EventBusName: aws.String(EventBusName),
			Detail:       aws.String(details),
			DetailType:   aws.String(EventDetail),
			Source:       aws.String(EventSource),
		}},
	}
}

func dispatchEvent(req events.ALBTargetGroupRequest, cfg aws.Config) error {
	srv := eventBridgeSession(cfg)
	details := string([]byte(req.Body))

	e := eventRequestEntry(details)
	request := srv.PutEventsRequest(&e)

	_, err := request.Send(context.TODO())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %s\n", err)
		return err
	}
	return nil
}

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
		return "", err
	}
	//As resp struct(obj) --> Nested struct parameter --> Value == *String
	return *resp.Parameter.Value, nil
}


func handle(req events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	// Parse needed values from GitHub webhook payload
	cfg := defaultConfig()

	secret, err := getSSMParameter(cfg)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %s\n", err)
		return events.ALBTargetGroupResponse{Body: err.Error(), StatusCode: 503}, nil
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

	err = dispatchEvent(req, cfg)
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



