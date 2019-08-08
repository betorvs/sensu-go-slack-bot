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
	outputParsed = "parsed"
	notFound     = "Not Found"
)

// sensuToken struct
type sensuToken struct {
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresAt    string `json:"expires_at,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// BasicAuth func go to sensu api and get a new token
func BasicAuth() (string, error) {
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
