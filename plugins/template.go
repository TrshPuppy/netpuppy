package plugins

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func init() {
    Register("template", &Template{})
}

type Template struct{}

func (r *Template) Execute(pluginDataChan chan<- string) {
    fmt.Println("[Plugin] This is a basic boilerplate template for a netpuppy plugin. type exit to quit.")
    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("> ")
        scanner.Scan()
        input := scanner.Text()

        if strings.TrimSpace(input) == "exit" {
            fmt.Println("[Plugin] goodbye!")
            break
        }

        pluginInput := "[Plugin] " + input

        // sends the input to the plugin channel instead of the main input channel.
        pluginDataChan <- pluginInput

        fmt.Println(pluginInput)
    }
}