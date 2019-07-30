package controller

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/betorvs/sensu-go-slack-bot/config"
	"github.com/betorvs/sensu-go-slack-bot/usecase"
	"github.com/labstack/echo"
	"github.com/nlopes/slack"
)

// ReceiveEvents func
func ReceiveEvents(c echo.Context) (err error) {
	// Thanks for https://medium.com/@xoen/golang-read-from-an-io-readwriter-without-loosing-its-content-2c6911805361
	// Read the content
	var bodyBytes []byte
	if c.Request().Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
	}
	// Restore the io.ReadCloser to its original state
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	// Use the content
	bodyString := string(bodyBytes)
	// Copy Forms to Struct
	data := new(slack.SlashCommand)
	data.Token = c.FormValue("token")
	data.TeamID = c.FormValue("team_id")
	data.TeamDomain = c.FormValue("team_domain")
	data.EnterpriseID = c.FormValue("enterprise_id")
	data.EnterpriseName = c.FormValue("enterprise_name")
	data.ChannelID = c.FormValue("channel_id")
	data.ChannelName = c.FormValue("channel_name")
	data.UserID = c.FormValue("user_id")
	data.UserName = c.FormValue("user_name")
	data.Command = c.FormValue("command")
	data.Text = c.FormValue("text")
	data.ResponseURL = c.FormValue("response_url")
	data.TriggerID = c.FormValue("trigger_id")
	// Headers
	slackRequestTimestamp := c.Request().Header.Get("X-Slack-Request-Timestamp")
	slackSignature := c.Request().Header.Get("X-Slack-Signature")

	basestring := fmt.Sprintf("v0:%s:%s", slackRequestTimestamp, bodyString)
	verifier := usecase.ValidateBot(slackSignature, basestring, config.SlackSigningSecret)
	if verifier != true {
		return c.JSON(http.StatusForbidden, nil)
	}
	go log.Printf("[AUDIT] User: %s, Executed: %s", data.UserName, data.Text)
	res, err := usecase.SlashCommandHandler(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusAccepted, res)
}
