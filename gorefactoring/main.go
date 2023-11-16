package main

import (
    "fmt"
    modules "main/modules"
)

const EXCHANGERATE_API_KEY = "51d9446da2379ccf84f73474"

func main() {
    data := modules.Currency(EXCHANGERATE_API_KEY)
    fmt.Print(data.String())
}