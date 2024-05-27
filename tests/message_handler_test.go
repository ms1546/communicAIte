package tests

import (
	"communicAIte/handlers"
	"testing"

	"github.com/slack-go/slack"
)

type mockRTM struct {
	sentMessages []string
}

func (m *mockRTM) SendMessage(msg *slack.OutgoingMessage) {
	m.sentMessages = append(m.sentMessages, msg.Text)
}

func (m *mockRTM) NewOutgoingMessage(text, channel string, options ...slack.RTMsgOption) *slack.OutgoingMessage {
	return &slack.OutgoingMessage{Text: text, Channel: channel}
}

func (m *mockRTM) GetUserInfo(userID string) (*slack.User, error) {
	return &slack.User{ID: userID, RealName: "Test User"}, nil
}

func TestHandleMessageEvent(t *testing.T) {
	openAIAPIKey := "test-openai-api-key"
	slackChannelID := "test-channel-id"
	mockRTM := &mockRTM{}

	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{"hello", "hello", "Hello Test User! How can I assist you today?"},
		{"help", "help", "Here are some commands you can use:\n- `hello`: Greet the bot\n- `help`: Display this help message"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := &slack.MessageEvent{
				Msg: slack.Msg{
					Channel: slackChannelID,
					Text:    tt.message,
					User:    "test-user-id",
				},
			}
			handlers.HandleMessageEvent(ev, mockRTM, openAIAPIKey, slackChannelID)

			if len(mockRTM.sentMessages) == 0 {
				t.Fatalf("Expected a response message but got none")
			}

			if mockRTM.sentMessages[len(mockRTM.sentMessages)-1] != tt.expected {
				t.Errorf("Expected response message %q but got %q", tt.expected, mockRTM.sentMessages[len(mockRTM.sentMessages)-1])
			}
		})
	}
}
