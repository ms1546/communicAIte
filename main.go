package main

import (
	"log"

	"communicAIte/handlers"
	"communicAIte/secrets"

	"github.com/slack-go/slack"
)

func main() {
	secretName := "MyAppSecrets"
	region := "ap-northeast-1"

	secret, err := secrets.GetSecret(secretName, region)
	if err != nil {
		log.Fatalf("Failed to retrieve secrets: %v", err)
	}

	api := slack.New(secret.SlackBotToken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			handlers.HandleMessageEvent(ev, rtm, secret.OpenAIAPIKey, secret.SlackChannelID)
		}
	}
}
