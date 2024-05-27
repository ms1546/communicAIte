package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/slack-go/slack"
)

type SlackRTM interface {
	SendMessage(msg *slack.OutgoingMessage)
	NewOutgoingMessage(text, channel string, options ...slack.RTMsgOption) *slack.OutgoingMessage
	GetUserInfo(userID string) (*slack.User, error)
}

func HandleMessageEvent(ev *slack.MessageEvent, rtm SlackRTM, openAIAPIKey, slackChannelID string) {
	if ev.Channel != slackChannelID {
		return
	}

	userMessage := ev.Text
	user, err := rtm.GetUserInfo(ev.User)
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		return
	}

	log.Printf("Received message from %s: %s", user.RealName, userMessage)

	if userMessage == "hello" {
		responseMessage := fmt.Sprintf("Hello %s! How can I assist you today?", user.RealName)
		log.Printf("Sent response: %s", responseMessage)
		rtm.SendMessage(rtm.NewOutgoingMessage(responseMessage, ev.Channel))
		return
	}

	if userMessage == "help" {
		responseMessage := "Here are some commands you can use:\n- `hello`: Greet the bot\n- `help`: Display this help message"
		log.Printf("Sent response: %s", responseMessage)
		rtm.SendMessage(rtm.NewOutgoingMessage(responseMessage, ev.Channel))
		return
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+openAIAPIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":      "text-davinci-003",
			"prompt":     userMessage,
			"max_tokens": 150,
		}).
		Post("https://api.openai.com/v1/completions")

	if err != nil {
		log.Printf("Error making request to OpenAI: %v", err)
		rtm.SendMessage(rtm.NewOutgoingMessage("Sorry, something went wrong while contacting the OpenAI API.", ev.Channel))
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		log.Printf("Error parsing response from OpenAI: %v", err)
		rtm.SendMessage(rtm.NewOutgoingMessage("Sorry, something went wrong while processing the response from OpenAI.", ev.Channel))
		return
	}

	responseMessage, ok := result["choices"].([]interface{})[0].(map[string]interface{})["text"].(string)
	if !ok {
		log.Printf("Unexpected response format from OpenAI: %v", result)
		rtm.SendMessage(rtm.NewOutgoingMessage("Sorry, I received an unexpected response from OpenAI.", ev.Channel))
		return
	}

	log.Printf("Sent response: %s", responseMessage)

	rtm.SendMessage(rtm.NewOutgoingMessage(responseMessage, ev.Channel))
}
