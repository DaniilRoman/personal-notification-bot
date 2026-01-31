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

func Test_GetUpdateStrings(t *testing.T) {
	// Test with nil BlogUpdateData
	var nilData *BlogUpdateData
	assert.Empty(t, nilData.GetUpdateStrings())
	
	// Test with empty BlogUpdateData
	emptyData := &BlogUpdateData{Blogs: []blogUpdate{}}
	assert.Empty(t, emptyData.GetUpdateStrings())
	
	// Test with blog updates
	blogData := &BlogUpdateData{
		Blogs: []blogUpdate{
			{
				Title: "Test Blog 1",
				Link:  "https://example.com/blog1",
			},
			{
				Title: "Test Blog 2",
				Link:  "https://medium.com/test/blog2",
			},
		},
	}
	
	updateStrings := blogData.GetUpdateStrings()
	assert.Len(t, updateStrings, 2)
	assert.Equal(t, "- Test Blog 1 {[example.com](https://example.com/blog1)}", updateStrings[0])
	assert.Equal(t, "- Test Blog 2 {[test](https://medium.com/test/blog2)}", updateStrings[1])
}
