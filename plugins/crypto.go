package plugins

import (
    "bufio"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strings"
    "netpuppy/utils"
)

func init() {
    Register("crypto", &CryptoTicker{})
}

type CryptoTicker struct{}

func (ct *CryptoTicker) Execute(pluginDataChan chan<- string) {
    fmt.Println(utils.Color("[CryptoTicker]", utils.Yellow) + " Type any cryptocurrency ID (e.g., bitcoin, ethereum) to get the price. Type 'exit' to quit.")
    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("> ")
        scanner.Scan()
        input := strings.ToLower(strings.TrimSpace(scanner.Text()))

        if input == "exit" {
            fmt.Println(utils.Color("[CryptoTicker]", utils.Yellow) + " Goodbye!")
            break
        }

        // Fetch the price using CoinGecko API
        url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", input)
        resp, err := http.Get(url)
        if err != nil {
            fmt.Printf(utils.Color("[CryptoTicker]", utils.Yellow) + " Failed to get price for %s: %v\n", input, err)
            continue
        }
        defer resp.Body.Close()

        var result map[string]map[string]float64
        if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
            fmt.Printf(utils.Color("[CryptoTicker]", utils.Yellow) + " Error decoding response for %s: %v\n", input, err)
            continue
        }

        if priceInfo, exists := result[input]; exists {
            response := fmt.Sprintf(utils.Color("[CryptoTicker]", utils.Yellow) + " Current price of %s is "+utils.Color("$%.2f", utils.Green), strings.ToUpper(input), priceInfo["usd"])
            pluginDataChan <- response
            fmt.Println(response)
        } else {
            fmt.Printf(utils.Color("[CryptoTicker]", utils.Yellow) + " Price information for %s not found.\n", input)
        }
    }
}