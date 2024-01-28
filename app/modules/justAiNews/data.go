package word

type JustAiNewsData struct {
	data string
}

func (d *JustAiNewsData) String() string {
	if d == nil {
		return ""
	}
    return d.data
}
