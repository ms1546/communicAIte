package tests

import (
	"communicAIte/secrets"
	"testing"
)

func TestGetSecret(t *testing.T) {
	secretName := "MyAppSecrets"
	region := "ap-northeast-1"

	secret, err := secrets.GetSecret(secretName, region)
	if err != nil {
		t.Fatalf("Failed to retrieve secrets: %v", err)
	}

	if secret.OpenAIAPIKey == "" {
		t.Error("Expected OpenAI API key to be set, but it was empty")
	}
	if secret.SlackBotToken == "" {
		t.Error("Expected Slack bot token to be set, but it was empty")
	}
	if secret.SlackChannelID == "" {
		t.Error("Expected Slack channel ID to be set, but it was empty")
	}
}
