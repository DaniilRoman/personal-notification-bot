package currency

import "fmt"

type CurrencyData struct {
    KeyValues map[string]float32
}

func newCurrencyData() *CurrencyData {
    var data CurrencyData
    data.KeyValues = make(map[string]float32)
    return &data
}

func (c *CurrencyData) String() string {
    if c == nil {
		return ""
	}
    res := "Currencies:\n"
    for k, v := range c.KeyValues {
       res += fmt.Sprintf("%s: %.2f RUB\n", k, v)
    }
    return res
}
