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

       currencyData.keyValues[currency] = response.ConversionRates.Rub
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
    keyValues map[string]float32
}

func (c *CurrencyData) String() string {
    res := ""
    for k, v := range c.keyValues {
       res += fmt.Sprintf("%s: %f RUB\n", k, v)
    }
    return res
}

func newCurrencyData() *CurrencyData {
    var data CurrencyData
    data.keyValues = make(map[string]float32)
    return &data
}