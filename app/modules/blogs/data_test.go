package blogs

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_websiteName(t *testing.T) {
	actualUrl := websiteName("https://aws.amazon.com/ru/blogs/machine-learning/feed/")

	assert.Equal(t, "aws.amazon.com", actualUrl)
}

func Test_websiteName_medium(t *testing.T) {
	actualUrl := websiteName("https://medium.com/feed/paypal-tech")

	assert.Equal(t, "paypal-tech", actualUrl)
}

func Test_websiteName_habr(t *testing.T) {
	actualUrl := websiteName("https://habr.com/ru/rss/company/just_ai/blog/?fl=ru")

	assert.Equal(t, "just_ai", actualUrl)
}
