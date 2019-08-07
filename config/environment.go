package config

import (
	"log"
	"os"
)

var (
	// Port to be listen by application
	Port string
	// SlackToken string
	SlackToken string
	// SlackSigningSecret string
	SlackSigningSecret string
	// SlackChannel string
	SlackChannel string
	// SensuGoURL string
	SensuGoURL string
	// SensuGoUser string
	SensuGoUser string
	// SensuGoSecret string
	SensuGoSecret string
	// SensuSlashCommand string
	SensuSlashCommand string
)

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func init() {
	Port = getEnv("SERVER_PORT", "9090")
	SensuSlashCommand = getEnv("SENSU_SLASH_COMMAND", "/sensu-go")
	SlackSigningSecret = os.Getenv("SLACK_SIGNING_SECRET")
	if SlackSigningSecret == "" {
		log.Fatal("variable SLACK_SIGNING_SECRET not defined")
	}
	SlackToken = os.Getenv("SLACK_TOKEN")
	if SlackToken == "" {
		log.Fatal("variable SLACK_TOKEN not defined")
	}
	SlackChannel = os.Getenv("SLACK_CHANNEL")
	if SlackChannel == "" {
		log.Fatal("variable SLACK_CHANNEL not defined")
	}
	SensuGoURL = os.Getenv("SENSU_URL")
	if SensuGoURL == "" {
		log.Fatal("variable SENSU_URL not defined")
	}
	SensuGoUser = os.Getenv("SENSU_USER")
	if SensuGoUser == "" {
		log.Fatal("variable SENSU_USER not defined")
	}
	SensuGoSecret = os.Getenv("SENSU_SECRET")
	if SensuGoSecret == "" {
		log.Fatal("variable SENSU_SECRET not defined")
	}
}
