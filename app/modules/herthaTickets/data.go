package hertha

import "fmt"

type HerthaTicketsData struct {
	Data string
}

func (d *HerthaTicketsData) String() string {
	if d == nil || d.Data == "" {
		return ""
	}
	return fmt.Sprintf("[Hertha Berlin tickets](%s):\n%s", herthaTicketsUrl, d.Data)
}
