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
	"[Html page](https://daniilroman.github.io/personal-notification-bot/)"
}

func (b *blogUpdate) String() string {
	websiteName := strings.Split(strings.TrimPrefix(strings.TrimPrefix(b.Link, "https://"), "http://"), "/")[0]
	resArticleStr := fmt.Sprintf("- %s {[%s](%s)}", b.Title, websiteName, b.Link)
	return resArticleStr
}
