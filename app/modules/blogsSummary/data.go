package blogsSummary

type BlogsSummaryData struct {
    summary string
}

func (s *BlogsSummaryData) String() string {
    if s == nil {
		return ""
	}
	return s.summary
}