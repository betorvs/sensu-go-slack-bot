package test

import (
	"testing"

	"github.com/betorvs/sensu-go-slack-bot/usecase"
	"github.com/nlopes/slack"
)

func SlachCommandTest(t *testing.T) {
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

func ValidateBotTest(t *testing.T) {
	signing := "v0=f2db95cc2f422dd8bb86185ec34270df93c0af1c6f6f755622b90c941dda2d72"
	message := "Teste123"
	slackSigningSecret := "test123###"
	verifier := usecase.ValidateBot(signing, message, slackSigningSecret)
	if verifier != true {
		t.Fatalf("Error: ValidateBot usecase inst work properly")
	}
	signing2 := "v0=f2db95cc2f422dd8bb86185ec34270df93c0af1c6f6f755622b90c941dda2d72"
	message2 := "Teste"
	slackSigningSecret2 := "test123###"
	verifier2 := usecase.ValidateBot(signing2, message2, slackSigningSecret2)
	if verifier2 != false {
		t.Fatalf("Error: ValidateBot usecase inst work properly")
	}
}
