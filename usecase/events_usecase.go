package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/betorvs/sensu-go-slack-bot/config"
	"github.com/betorvs/sensu-go-slack-bot/gateway/sensuclient"
	"github.com/labstack/echo"
	"github.com/nlopes/slack"
)

// SlashCommandHandler func handle with event received and return a message to controller
func SlashCommandHandler(data *slack.SlashCommand, c echo.Context) (slack.Msg, error) {
	var res slack.Msg
	switch data.Command {
	case "/sensu-go":
		values := strings.Fields(data.Text)
		action := values[0]
		check := values[1]
		server := values[2]
		namespace := values[3]
		go sensuclient.Connect(action, check, server, namespace, data.UserID, data.ChannelID)
		text := fmt.Sprintf("Check: %s, Server: %s, Namespace: %s, Processing...", check, server, namespace)
		message := slack.Msg{
			ResponseType: "in_channel",
			Text:         text}
		res = message

	default:
		log.Printf("[ERROR] Invalid slash command received: %s", data.Command)
		message := slack.Msg{
			Text: "Invalid Slash Command"}
		res = message

	}
	return res, nil
}

// ValidateBot func to validate auth from Slack Bot
func ValidateBot(timestamp string, signing string, message string) bool {
	mac := hmac.New(sha256.New, []byte(config.SlackSigningSecret))
	if _, err := mac.Write([]byte(message)); err != nil {
		log.Printf("mac.Write(%v) failed\n", message)
		return false
	}
	calculatedMAC := "v0=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(calculatedMAC), []byte(signing))
}
