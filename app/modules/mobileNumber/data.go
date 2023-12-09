package mobilenumber


type MobileNumberData struct {
	data string
}

func (d *MobileNumberData) String() string {
	if d == nil {
		return ""
	}
    return d.data
}
