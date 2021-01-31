package potboiler

import (
	"bytes"
	"github.com/mainak90/git-webhook-to-eventbridge"
	"text/template"
)

// MessageTemplateData contains data sent to Slack
type MessageTemplateData struct {
	URL      string
	FullName string
	Version  string
	Notes    string
	Channel  string
}

// MessageResponse contains data received from Slack
type MessageResponse struct {
	OK bool `json:"ok"`
}

func messageFromRequest(request main.Request) ([]byte, error) {
	template := messageTemplateFromPayloadForChannel(request.Payload, request.Channel)

	return messageFromTemplate(template)
}

func messageTemplateFromPayloadForChannel(payload main.Payload, channel string) MessageTemplateData {
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