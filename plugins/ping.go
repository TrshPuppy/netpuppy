package plugins

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
    "time"
)

func init() {
    Register("ping", &ping{})
}

type ping struct{}

func (tp *ping) Execute(pluginDataChan chan<- string) {
    fmt.Println("[ping] Type ip:port to ping. Type 'exit' to quit.")
    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("> ")
        scanner.Scan()
        input := scanner.Text()

        if strings.TrimSpace(input) == "exit" {
            fmt.Println("[ping] Goodbye!")
            break
        }

        // Validate input format
        if !strings.Contains(input, ":") {
            fmt.Println("[ping Error] Invalid format. Please use host:port.")
            continue
        }

        start := time.Now()
        conn, err := net.DialTimeout("tcp", input, 5*time.Second)
        if err != nil {
            fmt.Printf("[ping Error] Could not connect to %s: %v\n", input, err)
            continue
        }
        conn.Close()
        elapsed := time.Since(start)

        pingResult := fmt.Sprintf("[ping] to %s took %s.", input, elapsed)
        pluginDataChan <- pingResult
        fmt.Println(pingResult)
    }
}
