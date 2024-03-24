package blogsStats

type BlogsStatsData struct {
    sundayStats string
    newWeekWords string
    monthStats string
}

func (c *BlogsStatsData) String() string {
    if c == nil {
		return ""
	}
    res := ""
    if c.sundayStats != "" {
        res += "Top 10 popular keywords for the week: " + c.sundayStats
    }
    if c.monthStats != "" {
        res += "\nTop 10 popular keywords for the month: " + c.monthStats
    }
    if c.monthStats != "" {
        res += "\nNew keywords of the month: " + c.newWeekWords
    }
    return res
}
