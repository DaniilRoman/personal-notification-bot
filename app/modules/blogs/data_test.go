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
	actualUrl := websiteName("https://medium.com/adevinta-tech-blog/oops-i-forgot-to-publish-how-can-i-connect-to-the-container-then-8391a3b76c71?source=rss----19a122f075bd---4")

	assert.Equal(t, "adevinta-tech-blog", actualUrl)
}

func Test_websiteName_habr(t *testing.T) {
	actualUrl := websiteName("https://habr.com/ru/companies/ozontech/articles/817737/?utm_source=habrahabr&utm_medium=rss&utm_campaign=corporate_blog")

	assert.Equal(t, "ozontech", actualUrl)
}
