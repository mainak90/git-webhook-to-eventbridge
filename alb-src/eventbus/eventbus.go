package eventbus

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"os"
)

var (
	EventBusName = os.Getenv("EVENT_BUS_NAME")
	EventSource = "gitwebhook.lambda"
	EventDetail = os.Getenv("EVENT_DETAIL")
)

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

func DispatchEvent(req events.ALBTargetGroupRequest, cfg aws.Config) error {
	srv := eventBridgeSession(cfg)
	details := string([]byte(req.Body))

	e := eventRequestEntry(details)
	request := srv.PutEventsRequest(&e)

	_, err := request.Send(context.TODO())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %s\n", err)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eventbridge.ErrCodeConcurrentModificationException:
				return fmt.Errorf(eventbridge.ErrCodeConcurrentModificationException)
			case eventbridge.ErrCodeInternalException:
				return fmt.Errorf(eventbridge.ErrCodeInternalException)
			case eventbridge.ErrCodeOperationDisabledException:
				return fmt.Errorf(eventbridge.ErrCodeOperationDisabledException)
			case eventbridge.ErrCodeResourceNotFoundException:
				return fmt.Errorf(eventbridge.ErrCodeResourceNotFoundException)
			}
		} else {
			return fmt.Errorf(err.Error())
		}
	}
	return nil
}
