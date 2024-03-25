package plugins

import (
    "bufio"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
	"netpuppy/utils"
)


var RepositoryURL = "https://raw.githubusercontent.com/pianoplayerjames/netpuppy_plugins/main/"

func init() {
    Register("install", &install{})
}

// Command interface that all commands must implement.
type Command interface {
    Execute(pluginDataChan chan<- string)
}

// Global registry for commands.
var Commands = make(map[string]Command)

// Register function to add commands to the registry.
func Register(name string, cmd Command) {
    if _, exists := Commands[name]; exists {
        fmt.Printf("Warning: Command '%s' is already registered and will be overwritten.\n", name)
    }
    Commands[name] = cmd
}

type install struct{}

func (r *install) Execute(pluginDataChan chan<- string) {
    fmt.Println("[Plugin] Type the name of the plugin to install or 'exit' to quit.")
    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("> ")
        if !scanner.Scan() {
            return
        }
        input := strings.TrimSpace(scanner.Text())

        if input == "exit" {
            fmt.Println("[Plugin] Goodbye!")
            return
        }

        fileName := input
        if !strings.HasSuffix(fileName, ".go") {
            fileName += ".go"
        }

        rawURL := fmt.Sprintf("%s%s", RepositoryURL, fileName)

        if err := downloadFile(fmt.Sprintf("./plugins/%s", fileName), rawURL); err != nil {
            fmt.Printf(utils.Color("[Plugin] Error downloading file: %v\n", utils.Red), err)
            continue
        }

        fmt.Printf(utils.Color("[Plugin] '%s' installed successfully. Please restart netpuppy.\n", utils.Green), fileName)

        runCommand := strings.TrimSuffix(fileName, ".go")
        fmt.Printf(utils.Color("You can run the plugin by typing: '%s'\n", utils.Yellow), runCommand)

        break
    }
}

func downloadFile(filepath, url string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Check the status code
    if resp.StatusCode == http.StatusBadRequest {
        return fmt.Errorf(utils.Color("received a 400 Bad Request error for URL: %s", utils.Red), url)
    } else if resp.StatusCode >= 400 {
        return fmt.Errorf(utils.Color("received an HTTP error: %d - %s for URL: %s", utils.Red), resp.StatusCode, resp.Status, url)
    }

    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    return err
}
