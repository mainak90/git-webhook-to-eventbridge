package src

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
	"os"
)

const (
	EventBusName = "auto1-central"
	EventSource = "gitwebhook.lambda"
	EventDetail = "github-webhook"
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


func handle(req events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	// Parse needed values from GitHub webhook payload
	cfg := defaultConfig()

	var message interface{}
	err := json.Unmarshal([]byte(req.Body), &message)
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



