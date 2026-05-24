package blogssummary

import (
	"fmt"
	"log"
	"strings"

	"main/modules/blogs/data"
	"main/utils"
)

func GenerateAndSave(blogsData *data.BlogUpdateData, chatGpt *utils.ChatGptService, db *utils.DynamoDbService) {
	if blogsData == nil || len(blogsData.Blogs) == 0 {
		log.Println("No blog posts to summarize")
		return
	}

	var parts []string
	for _, blog := range blogsData.Blogs {
		desc := blog.Description
		if desc == "" {
			desc = blog.Title
		}
		parts = append(parts, fmt.Sprintf("Title: %s\nDescription: %s\nURL: %s", blog.Title, desc, blog.Link))
	}
	input := strings.Join(parts, "\n\n")

	summary := chatGpt.BlogsPodcastSummary(input)
	if summary == "" {
		log.Println("Empty podcast summary returned, skipping save")
		return
	}

	db.SaveItem("blogs_summary", summary)
	log.Println("Saved blogs_summary to DynamoDB")
}
