package utils

const (
    Black   = "\033[30m"
    Red     = "\033[31m"
    Green   = "\033[32m"
    Yellow  = "\033[33m"
    Blue    = "\033[34m"
    Magenta = "\033[35m"
    Cyan    = "\033[36m"
    White   = "\033[37m"
    Reset   = "\033[0m"
)

func Color(text string, color string) string {
	return color + text + Reset
}
