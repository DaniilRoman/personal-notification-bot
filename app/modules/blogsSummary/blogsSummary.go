package blogsSummary

import (
	"main/modules/blogs"
	"main/utils"
	"strings"
)

func BlogsSummary(blogsData *blogs.BlogUpdateData, chatGpt *utils.ChatGptService) *BlogsSummaryData {
	if blogsData == nil {
		return nil
	}
	elements := []string{}
	for _, blog := range blogsData.Blogs {
		elements = append(elements, blog.Summary)
	}
	allPartialSummary := strings.Join(elements, "\n")
	commonSummary := chatGpt.CommonSummaryFromSeveralNews(allPartialSummary)
	return &BlogsSummaryData{commonSummary}
}