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
)

const (
	outputJSON   = "json"
	outputEvent  = "event"
	outputCheck  = "check"
	outputEntity = "entity"
	notFound     = "Not Found"
)

// SensuToken struct
type SensuToken struct {
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresAt    string `json:"expires_at,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

var currentToken = SensuToken{}

// BasicAuth func go to sensu api and get a new token
func BasicAuth() (*SensuToken, error) {
	if currentToken.AccessToken != "" &&
		currentToken.ExpiresAt > string(time.Now().Unix()) {
		return &currentToken, nil
	}

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
	currentToken := new(SensuToken)
	json.Unmarshal(bodyText, &currentToken)
	defer resp.Body.Close()
	return currentToken, nil
}

// SensuGet func
func SensuGet(token string, url string, output string) (string, string, error) {
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
	case outputEvent:
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

	case outputEntity:
		if resp.StatusCode == 200 {
			var results []map[string]interface{}
			json.Unmarshal(bodyText, &results)
			for _, result := range results {
				entityClass := result["entity_class"]
				system := result["system"].(map[string]interface{})
				s += fmt.Sprintf("Hostname: %s, OS: %s %s, Version: %s, Entity Class: %s \n", system["hostname"], system["os"], system["platform"], system["platform_version"], entityClass)
			}
		} else {
			s = notFound
		}

	case outputCheck:
		if resp.StatusCode == 200 {
			var results []map[string]interface{}
			json.Unmarshal(bodyText, &results)
			for _, result := range results {
				command := result["command"]
				handlers := result["handlers"]
				subscriptions := result["subscriptions"]
				metadata := result["metadata"].(map[string]interface{})
				s += fmt.Sprintf("Check: %s, Command: %s, Namespace: %s, Subscriptions: %s, Handler: %s\n", metadata["name"], command, metadata["namespace"], subscriptions, handlers)
			}
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

// SensuPost func
func SensuPost(token string, url string, body []byte) (string, string, error) {
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
