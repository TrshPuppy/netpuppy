package plugins

import "fmt"

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