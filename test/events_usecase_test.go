package test

import (
	"testing"

	"github.com/betorvs/sensu-go-slack-bot/usecase"
	"github.com/nlopes/slack"
)

func TestRun(t *testing.T) {
	data := new(slack.SlashCommand)
	data.Command = "/sensu-go"
	data.Text = "get list-process server prod prod"
	response1 := "Not processed: Wrong number of fields 5"
	msg, _ := usecase.SlashCommandHandler(data)
	if msg.Text != response1 {
		t.Fatalf("Error: Expected 5 fields")
	}
	data2 := new(slack.SlashCommand)
	data2.Command = "/sensu-go"
	data2.Text = "get list-process server"
	response2 := "Not processed: Wrong number of fields 3"
	msg2, _ := usecase.SlashCommandHandler(data2)
	if msg2.Text != response2 {
		t.Fatalf("Error: Expected 3 fields")
	}
}
