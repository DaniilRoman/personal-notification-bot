package word

type WordOfTheDayData struct {
	data string
}

func (d *WordOfTheDayData) String() string {
	if d == nil {
		return ""
	}
    return d.data
}
