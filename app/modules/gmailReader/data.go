package gmailreader

import (
	"fmt"
	"strings"
)

type GmailReaderData struct {
	Messages []GmailMessage
}

type GmailMessage struct {
	Subject   string
	From      string
	Body      string
	MessageId string
}

func NewGmailMessage(subject, from, body, messageId string) GmailMessage {
	return GmailMessage{subject, from, body, messageId}
}

func (d *GmailReaderData) String() string {
	if d == nil || len(d.Messages) == 0 {
		return ""
	}
	msgStrings := make([]string, len(d.Messages))
	for i, msg := range d.Messages {
		msgStrings[i] = msg.String()
	}
	return "Gmail updates:\n\n" + strings.Join(msgStrings, "\n\n")
}

func (d *GmailReaderData) GetUpdateStrings() []string {
	if d == nil || len(d.Messages) == 0 {
		return []string{}
	}
	updateStrings := make([]string, len(d.Messages))
	for i, msg := range d.Messages {
		updateStrings[i] = msg.String()
	}
	return updateStrings
}

func (m *GmailMessage) String() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("From: %s", m.From))
	parts = append(parts, fmt.Sprintf("Subject: %s", m.Subject))

	if m.Body != "" {
		parts = append(parts, "")
		parts = append(parts, m.Body)
	}

	parts = append(parts, "")
	parts = append(parts, fmt.Sprintf("[Open in Gmail](%s)", m.MessageId))

	return strings.Join(parts, "\n")
}
