package slackclient

import (
	"fmt"
	"log"
	"net/http"
	"time"

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

// HealthSlack func
func HealthSlack(url string, token string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("[ERROR]: %s", err)
		return "", err
	}
	var bearer = "Bearer " + token
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR]: %s", err)
		return "", err
	}
	defer resp.Body.Close()
	return resp.Status, nil
}
