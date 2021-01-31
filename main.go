package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"bytes"
	"text/template"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"time"
)

type MessageTemplateData struct {
	URL      string
	FullName string
	Version  string
	Notes    string
	Channel  string
}

type Request struct {
	Channel string
	Payload Payload
}

// MessageResponse contains data received from Slack
type MessageResponse struct {
	OK bool `json:"ok"`
}

// Payload contains data received from GitHub
type Payload struct {
	Release    PayloadRelease    `json:"release"`
	Repository PayloadRepository `json:"repository"`
	Sender     PayloadSender     `json:"sender"`
}

// PayloadSender contains data received from GitHub about the user
type PayloadSender struct {
	Name string `json:"login"`
}

// PayloadRepository contains data received from GitHub about the repository
type PayloadRepository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	URL      string `json:"html_url"`
}

// PayloadRelease contains data received from GitHub about the release
type PayloadRelease struct {
	Author     PayloadReleaseAuthor `json:"author"`
	Name       string               `json:"name"`
	Body       string               `json:"body"`
	Date       time.Time            `json:"created_at"`
	Draft      bool                 `json:"draft"`
	Prerelease bool                 `json:"prerelease"`
}

// PayloadReleaseAuthor contains data received from GitHub about the release author
type PayloadReleaseAuthor struct {
	Name string `json:"login"`
	URL  string `json:"html_url"`
}

func messageFromRequest(request Request) ([]byte, error) {
	template := messageTemplateFromPayloadForChannel(request.Payload, request.Channel)

	return messageFromTemplate(template)
}

func messageTemplateFromPayloadForChannel(payload Payload, channel string) MessageTemplateData {
	return MessageTemplateData{
		payload.Repository.URL,
		payload.Repository.FullName,
		payload.Release.Name,
		payload.Release.Body,
		channel,
	}
}

func messageFromTemplate(data MessageTemplateData) ([]byte, error) {
	t, _ := template.ParseFiles("templates/message.json")

	var message bytes.Buffer
	if err := t.Execute(&message, data); err != nil {
		return nil, err
	}

	return message.Bytes(), nil
}

func parseRequest(req events.APIGatewayProxyRequest) (*Request, error) {
	payload, err := parsePayload([]byte(req.Body))

	if err != nil {
		return nil, errors.New("Unable to parse request payload")
	}

	return &Request{
		req.PathParameters["channel"],
		*payload,
	}, nil
}

func parsePayload(data []byte) (*Payload, error) {
	var payload Payload

	err := json.Unmarshal(data, &payload)

	if err != nil {
		return nil, err
	}

	return &payload, nil
}

func handle(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse needed values from GitHub webhook payload
	request, err := parseRequest(req)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Unable to handle request", StatusCode: 500}, nil
	}

	// Create message from request
	message, err := messageFromRequest(*request)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Unable to create message", StatusCode: 500}, nil
	}

	fmt.Println(message)

	// Send response
	return events.APIGatewayProxyResponse{Body: "{ \"done\": true }", StatusCode: 200}, nil
}

func main() {
	lambda.Start(handle)
}


