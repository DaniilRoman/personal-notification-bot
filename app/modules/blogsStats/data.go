package blogsStats

type BlogsStatsData struct {
    data string
}

func (c *BlogsStatsData) String() string {
    if c == nil {
		return ""
	}
    return "Top 10 popular keywords for the week: " + c.data
}
