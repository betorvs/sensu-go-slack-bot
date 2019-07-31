package controller

import (
	"fmt"
	"net/http"

	"github.com/betorvs/sensu-go-slack-bot/config"
	"github.com/betorvs/sensu-go-slack-bot/gateway/sensuclient"
	"github.com/betorvs/sensu-go-slack-bot/gateway/slackclient"
	"github.com/labstack/echo"
)

// Health struct
type Health struct {
	Status string `json:"status"`
}

// CompleteHealth struct
type CompleteHealth struct {
	Status   string `json:"status"`
	SensuAPI string `json:"sensuapi"`
	SlackAPI string `json:"slackapi"`
}

// CheckHealth func to return OK and http 200
func CheckHealth(c echo.Context) error {
	health := Health{}
	health.Status = "UP"
	return c.JSON(http.StatusOK, health)
}

// CompleteCheck func is used to check connections between this bot to SensuGo and Slack.
// Don't use it as a test probe in Kubernetes
func CompleteCheck(c echo.Context) error {
	completeHealth := CompleteHealth{}
	completeHealth.Status = "OK"
	sensuURL := fmt.Sprintf("%s/health", config.SensuGoURL)
	sensu, err := sensuclient.SensuHealth(sensuURL)
	if err != nil {
		completeHealth.SensuAPI = "NOT OK"
	} else {
		completeHealth.SensuAPI = sensu
	}
	slackURL := "https://slack.com/api/auth.test"
	slack, err := slackclient.HealthSlack(slackURL, config.SlackToken)
	if err != nil {
		completeHealth.SlackAPI = "NOT OK"
	} else {
		completeHealth.SlackAPI = slack
	}
	return c.JSON(http.StatusOK, completeHealth)
}
