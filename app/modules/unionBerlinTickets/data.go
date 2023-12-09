package union

import "fmt"

type UnionBerlinTicketsData struct {
	data string
}

func (d *UnionBerlinTicketsData) String() string {
	if d == nil || d.data == "" {
		return ""
	}
	return fmt.Sprintf("[Union Berlin tickets](%s):\n%s", unionBerlinUrl, d.data)
}
