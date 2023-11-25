package modules

import (
    "fmt"
    "log"
    utils "main/utils"
)

func Currency(token string) *CurrencyData {
    currencyData, err := currency(token)
    if err != nil {
       log.Printf("Error in currency module: %s", err)
    }
    return currencyData
}

func currency(token string) (*CurrencyData, error) {
    currencyData := newCurrencyData()
    currencies := [2]string{"USD", "EUR"}
    for _, currency := range currencies {
       url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", token, currency)
       response := currencyResponse{}
       err := utils.DoGet(url, &response)
       if err != nil {
          return nil, err
       }

       currencyData.KeyValues[currency] = response.ConversionRates.Rub
    }

    return currencyData, nil
}

type currencyResponse struct {
    ConversionRates currencyEntry `json:"conversion_rates"`
}

type currencyEntry struct {
    Rub float32 `json:"RUB"`
}

type CurrencyData struct {
    KeyValues map[string]float32
}

func (c *CurrencyData) String() string {
    if c == nil {
		return ""
	}
    res := "Currencies:\n"
    for k, v := range c.KeyValues {
       res += fmt.Sprintf("%s: %.2f RUB", k, v)
    }
    return res
}

func newCurrencyData() *CurrencyData {
    var data CurrencyData
    data.KeyValues = make(map[string]float32)
    return &data
}