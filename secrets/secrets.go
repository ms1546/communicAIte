package secrets

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type Secret struct {
	OpenAIAPIKey   string `json:"OPENAI_API_KEY"`
	SlackBotToken  string `json:"SLACK_BOT_TOKEN"`
	SlackChannelID string `json:"SLACK_CHANNEL_ID"`
}

func GetSecret(secretName, region string) (Secret, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return Secret{}, err
	}
	svc := secretsmanager.New(sess)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return Secret{}, err
	}

	if result.SecretString == nil {
		return Secret{}, errors.New("SecretString is nil")
	}

	var secret Secret
	if err := json.Unmarshal([]byte(*result.SecretString), &secret); err != nil {
		return Secret{}, err
	}

	return secret, nil
}
