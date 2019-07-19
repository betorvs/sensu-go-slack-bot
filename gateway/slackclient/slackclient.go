package slackclient

import (
	"fmt"
	"log"

	"github.com/betorvs/sensu-go-slack-bot/config"
	"github.com/nlopes/slack"
)

// Slack struct for our slackbot
type Slack struct {
	Name  string
	Token string

	User   string
	UserID string

	Client *slack.Client
}

// New returns a new instance of the Slack struct, primary for our slackbot
func New() (*Slack, error) {
	return &Slack{Client: slack.New(config.SlackToken), Token: config.SlackToken, Name: "Slack to Sensu Go"}, nil
}

// EphemeralMessage func
func EphemeralMessage(channel string, user string, message string) error {
	s, err := New()
	if err != nil {
		log.Printf("Error creating slack client: %s", err)
	}
	attachment := slack.Attachment{
		Text: message,
	}
	if _, err := s.Client.PostEphemeral(channel, user, slack.MsgOptionAttachments(attachment)); err != nil {
		return fmt.Errorf("failed to post message: %v", err)
	}
	return nil
}

// EphemeralFileMessage func
func EphemeralFileMessage(channel string, user string, message string, title string) error {
	s, err := New()
	if err != nil {
		log.Printf("Error creating slack client: %s", err)
	}
	if message == "" {
		message = "Silenced"
	}
	params := slack.FileUploadParameters{
		Filename: "result.txt", Title: title, Content: message,
		Channels: []string{channel}}
	if _, err := s.Client.UploadFile(params); err != nil {
		log.Printf("Unexpected error: %s", err)
	}
	return nil
}
