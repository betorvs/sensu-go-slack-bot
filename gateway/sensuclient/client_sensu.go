package sensuclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/betorvs/sensu-go-slack-bot/config"
	"github.com/betorvs/sensu-go-slack-bot/gateway/slackclient"
)

const (
	outputJSON   = "json"
	outputParsed = "parsed"
	notFound     = "Not Found"
)

// sensuToken struct
type sensuToken struct {
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresAt    string `json:"expires_at,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// payload struct for post in sensu
type payload struct {
	Check         string   `json:"check,omitempty"`
	Subscriptions []string `json:"subscriptions,omitempty"`
}

// silence struct
type silence struct {
	Metadata        metadata `json:"metadata"`
	Check           string   `json:"check"`
	Expire          int      `json:"expire"`
	ExpireOnResolve bool     `json:"expire_on_resolve"`
	Subscription    string   `json:"subscription"`
}

type metadata struct {
	Namespace string `json:"namespace"`
}

func basicAuth() (string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	sensuURL := fmt.Sprintf("%s/auth", config.SensuGoURL)
	req, err := http.NewRequest("GET", sensuURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(config.SensuGoUser, config.SensuGoSecret)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR]: %s", err)
	}
	data := new(sensuToken)
	json.Unmarshal(bodyText, &data)
	defer resp.Body.Close()
	return data.AccessToken, nil
}

// Connect func func Connect(action string, check string, server string, namespace string, userid string, channel string) (string, error)
func Connect(action string, check string, server string, namespace string, userid string, channel string) {
	token, err := basicAuth()
	if err != nil {
		log.Printf("[ERROR] Auth Problems: %s", err)
	}
	switch action {
	case "get":
		sensuURL := fmt.Sprintf("%s/api/core/v2/namespaces/%s/events/%s/%s", config.SensuGoURL, namespace, server, check)
		s, body, err := sensuGet(token, sensuURL, outputParsed)
		if err != nil {
			log.Printf("[ERROR]: %s", err)
		}
		go slackclient.EphemeralFileMessage(channel, userid, body, s)

	case "execute":
		entity := fmt.Sprintf("entity:%s", server)
		formPost := payload{
			Check:         check,
			Subscriptions: []string{entity},
		}
		bodymarshal, err := json.Marshal(&formPost)
		if err != nil {
			log.Printf("[ERROR] %s", err)
		}
		sensuURL := fmt.Sprintf("%s/api/core/v2/namespaces/%s/checks/%s/execute", config.SensuGoURL, namespace, check)
		s, _, err := sensuPost(token, sensuURL, bodymarshal)
		if err != nil {
			log.Printf("[ERROR]: %s", err)
		}
		text := fmt.Sprintf("Check: %s, Server: %s, Namespace: %s, Response: %s", check, server, namespace, s)
		go slackclient.EphemeralMessage(channel, userid, text)

	case "silence":
		entity := fmt.Sprintf("entity:%s", server)
		metadata := metadata{
			Namespace: namespace,
		}
		formPost := silence{
			Metadata:        metadata,
			Check:           check,
			Expire:          -1,
			ExpireOnResolve: true,
			Subscription:    entity,
		}
		bodymarshal, err := json.Marshal(&formPost)
		if err != nil {
			log.Printf("[ERROR] %s", err)
		}
		sensuURL := fmt.Sprintf("%s/api/core/v2/namespaces/%s/silenced", config.SensuGoURL, namespace)
		s, _, err := sensuPost(token, sensuURL, bodymarshal)
		if err != nil {
			log.Printf("[ERROR]: %s", err)
		}
		text := fmt.Sprintf("Check: %s, Server: %s, Namespace: %s, Response: %s", check, server, namespace, s)
		go slackclient.EphemeralMessage(channel, userid, text)

	case "describe":
		switch check {
		case "check":
			sensuURL := fmt.Sprintf("%s/api/core/v2/namespaces/%s/checks/%s", config.SensuGoURL, namespace, server)
			s, body, err := sensuGet(token, sensuURL, outputJSON)
			if err != nil {
				log.Printf("[ERROR]: %s", err)
			}
			go slackclient.EphemeralFileMessage(channel, userid, body, s)

		case "entity":
			sensuURL := fmt.Sprintf("%s/api/core/v2/namespaces/%s/entities/%s", config.SensuGoURL, namespace, server)
			s, body, err := sensuGet(token, sensuURL, outputJSON)
			if err != nil {
				log.Printf("[ERROR]: %s", err)
			}
			go slackclient.EphemeralFileMessage(channel, userid, body, s)

		default:
			log.Println("[ERROR] describe unknown field")
			s := fmt.Sprintf("Please Use: %s describe [check|entity] [name] [namespace]", config.SensuSlashCommand)
			go slackclient.EphemeralMessage(channel, userid, s)
		}

	default:
		log.Println("[ERROR] unknown action")
		s := fmt.Sprintf("Please Use: %s [get|execute|silence|describe] [check|entity] [name] [namespace]", config.SensuSlashCommand)
		go slackclient.EphemeralMessage(channel, userid, s)
	}
}

// sensuGet func
func sensuGet(token string, url string, output string) (string, string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("[ERROR]: %s", err)
		// log.Fatal(err)
	}
	var bearer = "Bearer " + token
	req.Header.Add("Authorization", bearer)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[ERROR]: %s", err)
		// log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] %s", err)
	}
	var s string
	switch output {
	case outputParsed:
		if resp.StatusCode == 200 {
			var result map[string]interface{}
			json.Unmarshal(bodyText, &result)
			check := result["check"].(map[string]interface{})
			entity := result["entity"].(map[string]interface{})
			details := entity["system"].(map[string]interface{})
			s = fmt.Sprintf("Hostname: %s, %s, %s, Check Output: \n%s", details["hostname"], details["platform"], details["platform_version"], check["output"])
		} else {
			s = notFound
		}

	case outputJSON:
		if resp.StatusCode == 200 {
			s = string(bodyText)
		} else {
			s = notFound
		}

	default:
		log.Println("[ERROR] unknown output method")
	}

	defer resp.Body.Close()
	return resp.Status, s, nil
}

// sensuPost func
func sensuPost(token string, url string, body []byte) (string, string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	var bearer = "Bearer " + token
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] %s", err)
	}
	s := string(bodyText)
	defer resp.Body.Close()
	return resp.Status, s, nil
}

// SensuHealth func
func SensuHealth(url string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 3,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("[ERROR]: %s", err)
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[ERROR]: %s", err)
		return "", err
	}
	defer resp.Body.Close()
	return resp.Status, nil
}
