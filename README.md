# Sensu-Go-Slack-bot

A simple slack bot create using slash commands integration to communicate with Sensu-Go API.

![](https://media.giphy.com/media/YSr3jO6YncW1FQ85Ll/giphy.gif)

## Create a App in Slack

Add feature Slash Command with these parameters:
* Command: `/sensu-go`
* Request URL: `https://URL/sensu-go-bot/v1/events`
* Short Description: `Talk with Monitoring System Sensu Go `
* Usage Hint: `[get|execute|silence] [check] [server] [namespace]`

TIP: instead namespace, configure with namespace itself, like: prod|stg.

Add feature "Bots" to this bot.


In Oauth permissions add:
* CONVERSATIONS chat:write:bot
* FILES files:write:user

Install these Application in a Channel.

## Create an user in Sensu-Go

I create the sensu-go-slack-bot using these commands:

```sh
sensuctl user create sensu-go-bot --password "LONGPASSWORD"

sensuctl cluster-role create sensu-go-bot-role --verb get,list,create,update --resource checks,events,silenced --namespace prod

sensuctl cluster-role-binding create sensu-go-bot-rolebinding --cluster-role=sensu-go-bot-role --user=sensu-go-bot --namespace prod

```

## Configuration 

You need to configure these for local tests or real deployment.

Configure these environment variables:
* **SENSU_USER**=sensu-go-bot : same user create in Sensu-Go API.
* **SENSU_URL**=https://SENSU-URL:8080 : Keep without a slash '/' in the end.
* **SENSU_SECRET**"" : Sensu bot password 
* **SLACK_TOKEN** : App token from Oauth in Slack (start with xoxb-)
* **SLACK_SIGNING_SECRET** : App Signing secret from Slack App.

In Kubernetes deployments you can use secrets for these 3 last variables.


## Docker Build 

```sh 
docker build -t betorvs/sensu-go-slack-bot:test1 -f Dockerfile .
```

## Deploy using Helm and Kubectl

A basic deploy:

```sh
kubect create ns sensu-go-slack-bot
kubectl apply -f deployment/sensu-go-slack-bot/secrets.yaml -n sensu-go-slack-bot
helm upgrade --install sensu-go-slack-bot deployment/sensu-go-slack-bot/ --namespace sensu-go-slack-bot
```

### In prod

Configure a proper values-prod.yaml and run:

```sh
helm upgrade --install sensu-go-slack-bot deployment/sensu-go-slack-bot/ -f sensu-go-slack-bot/values-prod.yaml --namespace sensu-go-slack-bot
```

### Private repositories

Include `ImagePullSecrets` in deployment/sensu-go-slack-bot/templates/deployment.yaml.


## Extras

Create these extras checks on Sensu-Go to create a possibility to run a small troubleshooting from these slash command:

```sh
sensuctl check create list-process --command 'ps -ef' --publish=false --interval 60 --subscriptions linux --handlers default --namespace default
```

To run on Slack (that channel where it was installed):

```
/sensu-go execute list-process server-hostname default
```

To get the results:
```
/sensu-go get list-process server-hostname default
```


## Reference

https://medium.com/@emilygoldfein/creating-slack-slash-commands-using-go-3cea3b3f0920

https://api.slack.com/docs/verifying-requests-from-slack

https://github.com/nlopes/slack

