package controller

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/betorvs/sensu-go-slack-bot/config"
	"github.com/betorvs/sensu-go-slack-bot/usecase"
	"github.com/labstack/echo"
	"github.com/nlopes/slack"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"code", "method"})
	histogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_requests_response_time",
		Help:    "Time take to answer",
		Buckets: []float64{1, 2, 5, 6, 10}, //defining small buckets as this app should not take more than 1 sec to respond
	}, []string{"code", "method"})
)

// registerMetrics func
func registerMetrics(code string, method string, start time.Duration) {
	prometheus.Register(httpRequestsTotal)
	prometheus.Register(histogram)
	histogram.WithLabelValues(code, method).Observe(start.Seconds())
	httpRequestsTotal.WithLabelValues(code, method).Inc()
}

// ReceiveEvents func
func ReceiveEvents(c echo.Context) (err error) {
	start := time.Now()
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
		go registerMetrics("403", "POST", time.Since(start))
		return c.JSON(http.StatusForbidden, nil)
	}
	if data.ChannelID != config.SlackChannel {
		go registerMetrics("403", "POST", time.Since(start))
		return c.JSON(http.StatusForbidden, nil)
	}
	go log.Printf("[AUDIT] User: %s, Executed: %s", data.UserName, data.Text)
	res, err := usecase.SlashCommandHandler(data)
	if err != nil {
		go registerMetrics("400", "POST", time.Since(start))
		return c.JSON(http.StatusBadRequest, err)
	}
	go registerMetrics("202", "POST", time.Since(start))
	return c.JSON(http.StatusAccepted, res)
}
