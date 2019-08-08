package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/betorvs/sensu-go-slack-bot/config"
	"github.com/nlopes/slack"
)

// SlashCommandHandler func handle with event received and return a message to controller
func SlashCommandHandler(data *slack.SlashCommand) (slack.Msg, error) {
	var res slack.Msg
	switch data.Command {
	case config.SensuSlashCommand:
		values := strings.Fields(data.Text)
		length := len(values)
		if length < 4 || length > 4 {
			text := fmt.Sprintf("Not processed: Wrong number of fields %s", strconv.Itoa(length))
			message := slack.Msg{
				ResponseType: "in_channel",
				Text:         text}
			res = message
		} else {
			action := values[0]
			check := values[1]
			server := values[2]
			namespace := values[3]
			go SensuConnect(action, check, server, namespace, data.UserID, data.ChannelID)
			text := fmt.Sprintf("Check: %s, Server: %s, Namespace: %s, Processing...", check, server, namespace)
			message := slack.Msg{
				ResponseType: "in_channel",
				Text:         text}
			res = message
		}

	default:
		log.Printf("[ERROR] Invalid slash command received: %s", data.Command)
		message := slack.Msg{
			Text: "Invalid Slash Command"}
		res = message

	}
	return res, nil
}

// ValidateBot func to validate auth from Slack Bot
func ValidateBot(signing string, message string, mysigning string) bool {
	mac := hmac.New(sha256.New, []byte(mysigning))
	if _, err := mac.Write([]byte(message)); err != nil {
		log.Printf("mac.Write(%v) failed\n", message)
		return false
	}
	calculatedMAC := "v0=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(calculatedMAC), []byte(signing))
}
