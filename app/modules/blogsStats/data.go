package blogsStats

type BlogsStatsData struct {
    sundayStats string
    monthStats string
    newWeekWords string
}

func (c *BlogsStatsData) String() string {
    if c == nil {
		return ""
	}
    res := ""
    if c.sundayStats != "" {
        res += "**Top 10 popular keywords for the week:** " + c.sundayStats
    }
    if c.monthStats != "" {
        res += "\n**Top 10 popular keywords for the month:** " + c.monthStats
    }
    if c.newWeekWords != "" {
        res += "\n**New keywords of the month:** " + c.newWeekWords
    }
    return res
}
