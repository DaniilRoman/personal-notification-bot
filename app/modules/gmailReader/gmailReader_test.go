package gmailreader

import (
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestGofeedParsesAtom30(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<feed version="0.3" xmlns="http://purl.org/atom/ns#">
  <title>Gmail - Inbox for test@gmail.com</title>
  <fullcount>2</fullcount>
  <entry>
    <title>Hello from GitHub</title>
    <summary>Your PR was merged</summary>
    <author>
      <name>GitHub</name>
      <email>noreply@github.com</email>
    </author>
    <link rel="alternate" href="https://mail.google.com/mail?message_id=123" type="text/html"/>
    <issued>2025-01-15T10:00:00Z</issued>
  </entry>
  <entry>
    <title>Build passed</title>
    <summary>CI pipeline succeeded</summary>
    <author>
      <name>CI Bot</name>
      <email>ci@build.com</email>
    </author>
    <link rel="alternate" href="https://mail.google.com/mail?message_id=456" type="text/html"/>
    <issued>2025-01-15T11:00:00Z</issued>
  </entry>
</feed>`

	parser := gofeed.NewParser()
	feed, err := parser.ParseString(xmlData)
	assert.NoError(t, err)
	assert.NotNil(t, feed)
	t.Logf("feed title: %q", feed.Title)
	t.Logf("items: %d", len(feed.Items))
	for i, item := range feed.Items {
		t.Logf("item %d: title=%q desc=%q authors=%v links=%v published=%v",
			i, item.Title, item.Description, item.Authors, item.Links, item.PublishedParsed)
	}
}

func TestMatchesSender(t *testing.T) {
	t.Run("exact match", func(t *testing.T) {
		assert.True(t, matchesSender("user@example.com", []string{"user@example.com"}))
	})

	t.Run("case insensitive", func(t *testing.T) {
		assert.True(t, matchesSender("User@Example.COM", []string{"user@example.com"}))
	})

	t.Run("partial match", func(t *testing.T) {
		assert.True(t, matchesSender("noreply@github.com", []string{"github.com"}))
	})

	t.Run("no match", func(t *testing.T) {
		assert.False(t, matchesSender("other@example.com", []string{"github.com"}))
	})

	t.Run("multiple senders", func(t *testing.T) {
		assert.True(t, matchesSender("bot@amazon.com", []string{"github.com", "amazon.com"}))
	})

	t.Run("empty senders list returns false", func(t *testing.T) {
		assert.False(t, matchesSender("test@test.com", []string{}))
	})
}

func TestGmailMessageString(t *testing.T) {
	t.Run("full message", func(t *testing.T) {
		msg := GmailMessage{
			Subject:   "Hello",
			From:      "Sender <sender@example.com>",
			Body:      "Body text",
			MessageId: "https://mail.google.com/mail?message_id=abc123",
		}
		result := msg.String()
		assert.Contains(t, result, "From: Sender <sender@example.com>")
		assert.Contains(t, result, "Subject: Hello")
		assert.Contains(t, result, "Body text")
		assert.Contains(t, result, "https://mail.google.com/mail?message_id=abc123")
	})

	t.Run("empty body", func(t *testing.T) {
		msg := GmailMessage{
			Subject:   "No Body",
			From:      "bot@test.com",
			Body:      "",
			MessageId: "https://mail.google.com/mail?message_id=msg456",
		}
		result := msg.String()
		assert.Contains(t, result, "Subject: No Body")
		assert.NotContains(t, result, "\n\n\n\n")
	})

	t.Run("gmail link format", func(t *testing.T) {
		msg := GmailMessage{
			Subject:   "Test",
			From:      "a@b.com",
			Body:      "body",
			MessageId: "https://mail.google.com/mail?message_id=msg_789",
		}
		result := msg.String()
		assert.Contains(t, result, "[Open in Gmail](https://mail.google.com/mail?message_id=msg_789)")
	})
}

func TestGmailReaderDataGetUpdateStrings(t *testing.T) {
	t.Run("nil data", func(t *testing.T) {
		var nilData *GmailReaderData
		assert.Empty(t, nilData.GetUpdateStrings())
	})

	t.Run("empty messages", func(t *testing.T) {
		data := &GmailReaderData{Messages: []GmailMessage{}}
		assert.Empty(t, data.GetUpdateStrings())
	})

	t.Run("multiple messages", func(t *testing.T) {
		data := &GmailReaderData{
			Messages: []GmailMessage{
				{Subject: "Subject 1", From: "a@b.com", Body: "Body 1", MessageId: "https://mail.google.com/mail?id=1"},
				{Subject: "Subject 2", From: "c@d.com", Body: "Body 2", MessageId: "https://mail.google.com/mail?id=2"},
			},
		}
		strings := data.GetUpdateStrings()
		assert.Len(t, strings, 2)
		assert.Contains(t, strings[0], "Subject 1")
		assert.Contains(t, strings[1], "Subject 2")
	})
}

func TestGmailReaderDataString(t *testing.T) {
	t.Run("nil data", func(t *testing.T) {
		var nilData *GmailReaderData
		assert.Empty(t, nilData.String())
	})

	t.Run("empty messages", func(t *testing.T) {
		data := &GmailReaderData{Messages: []GmailMessage{}}
		assert.Empty(t, data.String())
	})

	t.Run("multiple messages", func(t *testing.T) {
		data := &GmailReaderData{
			Messages: []GmailMessage{
				{Subject: "Subj", From: "a@b.com", Body: "", MessageId: "id1"},
			},
		}
		result := data.String()
		assert.Contains(t, result, "Gmail updates:")
		assert.Contains(t, result, "Subject: Subj")
	})
}


