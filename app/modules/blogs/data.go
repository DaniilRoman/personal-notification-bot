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
	PopularWords string
	Author string
}

func NewBlogUpdate(title string, link string, img string, summary string, popularWords string) blogUpdate {
	return blogUpdate{title, link, img, summary, popularWords, substractedUrl(link)}
}

func (blogs *BlogUpdateData) String() string {
	if blogs == nil || len(blogs.Blogs) == 0 {
		return ""
	}
	blogStrings := make([]string, len(blogs.Blogs))
	for i, blog := range blogs.Blogs {
		blogStrings[i] = blog.String()
	}
    return "Blogs updates:\n\n" + 
	strings.Join(blogStrings, "\n\n") + 
	"\n\n[Html page](https://daniilroman.github.io/personal-notification-bot/)"
}

func (blogs *BlogUpdateData) PopularWords() string {
	if blogs == nil {
		return ""
	}
	popularWords := make([]string, len(blogs.Blogs))
	for i, blog := range blogs.Blogs {
		popularWords[i] = blog.PopularWords
	}
	return strings.Join(popularWords, ",")
}

func (b *blogUpdate) String() string {
	websiteName := websiteName(b.Link)
	resArticleStr := fmt.Sprintf("- %s {[%s](%s)}", b.Title, websiteName, b.Link)
	return resArticleStr
}

func websiteName(link string) string {
	if strings.Contains(link, "medium.com") {
		urlParts := strings.Split(link, "/")
		if len(urlParts) > 3 {
		  return urlParts[3]
		}
	}
	if strings.Contains(link, "habr.com") {
		urlParts := strings.Split(link, "/")
		if len(urlParts) > 5 {
		  return urlParts[5]
		}
	}
	return substractedUrl(link)
}

func substractedUrl(link string) string {
	return  strings.Split(strings.TrimPrefix(strings.TrimPrefix(link, "https://"), "http://"), "/")[0]
}
