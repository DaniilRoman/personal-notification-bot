package currency

import (
    "fmt"
    "log"
    utils "main/utils"
)

const currencyUrl = "https://v6.exchangerate-api.com/v6/%s/latest/%s"

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
       url := fmt.Sprintf(currencyUrl, token, currency)
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
