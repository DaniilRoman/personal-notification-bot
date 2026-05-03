package gmailreader

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	gmailFeedURL       = "https://mail.google.com/mail/feed/atom/"
	gmailReadonlyScope = "https://www.googleapis.com/auth/gmail.readonly"
)

func GmailReader() *GmailReaderData {
	clientID := os.Getenv("GMAIL_CLIENT_ID")
	clientSecret := os.Getenv("GMAIL_CLIENT_SECRET")
	refreshToken := os.Getenv("GMAIL_REFRESH_TOKEN")
	senderEnv := os.Getenv("GMAIL_SENDER_EMAIL")

	if clientID == "" || clientSecret == "" || refreshToken == "" || senderEnv == "" {
		log.Printf("Gmail reader: missing required env vars (GMAIL_CLIENT_ID, GMAIL_CLIENT_SECRET, GMAIL_REFRESH_TOKEN, GMAIL_SENDER_EMAIL)")
		return nil
	}

	senders := strings.Split(senderEnv, ",")
	for i := range senders {
		senders[i] = strings.TrimSpace(senders[i])
	}

	res, err := gmailReader(clientID, clientSecret, refreshToken, senders)
	if err != nil {
		log.Printf("Error in Gmail reader module: %s", err)
	}
	return res
}

func gmailReader(clientID, clientSecret, refreshToken string, senders []string) (*GmailReaderData, error) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{gmailReadonlyScope},
	}

	token := &oauth2.Token{RefreshToken: refreshToken}
	client := config.Client(context.Background(), token)

	parser := gofeed.NewParser()
	parser.Client = client

	feed, err := parser.ParseURL(gmailFeedURL)
	if err != nil {
		return nil, fmt.Errorf("parsing feed: %w", err)
	}

	today := time.Now().Truncate(24 * time.Hour)

	var messages []GmailMessage
	for _, item := range feed.Items {
		if item.PublishedParsed == nil || item.PublishedParsed.Before(today) {
			continue
		}

		authorEmail := ""
		authorName := ""
		if len(item.Authors) > 0 {
			authorEmail = item.Authors[0].Email
			authorName = item.Authors[0].Name
		}
		if authorEmail == "" {
			continue
		}

		if !matchesSender(authorEmail, senders) {
			continue
		}

		link := ""
		if len(item.Links) > 0 {
			link = item.Links[0]
		}

		body := strings.TrimSpace(item.Description)
		runes := []rune(body)
		if len(runes) > 1000 {
			body = string(runes[:1000]) + "..."
		}

		messages = append(messages, GmailMessage{
			Subject:   item.Title,
			From:      fmt.Sprintf("%s <%s>", authorName, authorEmail),
			Body:      body,
			MessageId: link,
		})

		if len(messages) >= 2 {
			break
		}
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return &GmailReaderData{Messages: messages}, nil
}

func matchesSender(email string, senders []string) bool {
	lower := strings.ToLower(email)
	for _, s := range senders {
		if strings.Contains(lower, strings.ToLower(s)) {
			return true
		}
	}
	return false
}
