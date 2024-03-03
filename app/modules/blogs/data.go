package blogs

import (
	"fmt"
	"strings"
)

type BlogUpdateData struct {
	Blogs []blogUpdate 
}

type blogUpdate struct {
	Title string
	Link string
	Img string
	Summary string
	Author string
}

func NewBlogUpdate(title string, link string, img string, summary string) blogUpdate {
	return blogUpdate{title, link, img, summary, websiteName(link)}
}

func (c *BlogUpdateData) String() string {
	if c == nil {
		return ""
	}
	blogStrings := make([]string, len(c.Blogs))
	for i, blog := range c.Blogs {
		blogStrings[i] = blog.String()
	}
    return "Blogs updates:\n\n" + 
	strings.Join(blogStrings, "\n\n") + 
	"\n\n[Html page](https://daniilroman.github.io/personal-notification-bot/)"
}

func (b *blogUpdate) String() string {
	websiteName := websiteName(b.Link)
	resArticleStr := fmt.Sprintf("- %s {[%s](%s)}", b.Title, websiteName, b.Link)
	return resArticleStr
}

func websiteName(link string) string {
	return strings.Split(strings.TrimPrefix(strings.TrimPrefix(link, "https://"), "http://"), "/")[0]
}
