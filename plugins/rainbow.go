package plugins

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func init() {
    Register("rainbow", &Rainbow{})
}

type Rainbow struct{}

func (r *Rainbow) Execute(pluginDataChan chan<- string) {
    fmt.Println("[Rainbow] Type your message and see it turn into a rainbow!")
    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("> ")
        scanner.Scan()
        input := scanner.Text()

        if strings.TrimSpace(input) == "exit" {
            fmt.Println("Exiting rainbow plugin.")
            break
        }

        rainbowText := toRainbow(input)
        pluginDataChan <- rainbowText
        fmt.Println(rainbowText)
    }
}

func toRainbow(text string) string {
    rainbowColors := []string{"\033[31m", "\033[33m", "\033[32m", "\033[36m", "\033[34m", "\033[35m"}
    var rainbowText strings.Builder

    for i, char := range text {
        color := rainbowColors[i%len(rainbowColors)]
        rainbowText.WriteString(color + string(char))
    }
    rainbowText.WriteString("\033[0m")

    return rainbowText.String()
}