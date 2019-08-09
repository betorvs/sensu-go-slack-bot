package usecase

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/betorvs/sensu-go-slack-bot/config"
	"github.com/betorvs/sensu-go-slack-bot/gateway/sensuclient"
	"github.com/betorvs/sensu-go-slack-bot/gateway/slackclient"
)

// const
const (
	outputJSON   = "json"
	outputEvent  = "event"
	outputCheck  = "check"
	outputEntity = "entity"
	notFound     = "Not Found"
)

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

// metadata struct
type metadata struct {
	Namespace string `json:"namespace"`
}

// SensuConnect func func SensuConnect(action string, check string, server string, namespace string, userid string, channel string) (string, error)
func SensuConnect(action string, check string, server string, namespace string, userid string, channel string) {
	token, err := sensuclient.BasicAuth()
	if err != nil {
		log.Printf("[ERROR] Auth Problems: %s", err)
	}
	switch action {
	case "get":
		switch check {
		case "all":
			switch server {
			case "entity":
				sensuURL := fmt.Sprintf("%s/api/core/v2/namespaces/%s/entities", config.SensuGoURL, namespace)
				s, body, err := sensuclient.SensuGet(token, sensuURL, outputEntity)
				if err != nil {
					log.Printf("[ERROR]: %s", err)
				}
				go slackclient.EphemeralFileMessage(channel, userid, body, s)

			case "check":
				sensuURL := fmt.Sprintf("%s/api/core/v2/namespaces/%s/checks", config.SensuGoURL, namespace)
				s, body, err := sensuclient.SensuGet(token, sensuURL, outputCheck)
				if err != nil {
					log.Printf("[ERROR]: %s", err)
				}
				go slackclient.EphemeralFileMessage(channel, userid, body, s)

			default:
				log.Println("[ERROR] get all with 3rd field wrong")
				s := fmt.Sprintf("Please Use: %s get all [check|entity] [namespace]", config.SensuSlashCommand)
				go slackclient.EphemeralMessage(channel, userid, s)
			}

		default:
			sensuURL := fmt.Sprintf("%s/api/core/v2/namespaces/%s/events/%s/%s", config.SensuGoURL, namespace, server, check)
			s, body, err := sensuclient.SensuGet(token, sensuURL, outputEvent)
			if err != nil {
				log.Printf("[ERROR]: %s", err)
			}
			go slackclient.EphemeralFileMessage(channel, userid, body, s)
		}

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
		s, _, err := sensuclient.SensuPost(token, sensuURL, bodymarshal)
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
		s, _, err := sensuclient.SensuPost(token, sensuURL, bodymarshal)
		if err != nil {
			log.Printf("[ERROR]: %s", err)
		}
		text := fmt.Sprintf("Check: %s, Server: %s, Namespace: %s, Response: %s", check, server, namespace, s)
		go slackclient.EphemeralMessage(channel, userid, text)

	case "describe":
		switch check {
		case "check":
			sensuURL := fmt.Sprintf("%s/api/core/v2/namespaces/%s/checks/%s", config.SensuGoURL, namespace, server)
			s, body, err := sensuclient.SensuGet(token, sensuURL, outputJSON)
			if err != nil {
				log.Printf("[ERROR]: %s", err)
			}
			go slackclient.EphemeralFileMessage(channel, userid, body, s)

		case "entity":
			sensuURL := fmt.Sprintf("%s/api/core/v2/namespaces/%s/entities/%s", config.SensuGoURL, namespace, server)
			s, body, err := sensuclient.SensuGet(token, sensuURL, outputJSON)
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
