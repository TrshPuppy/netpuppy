package plugins

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"netpuppy/utils"
)

func init() {
	Register("animate", &Animate{})
}

type Animate struct{}

func (a *Animate) Execute(pluginDataChan chan<- string) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("Enter text to animate:")
    text, _ := reader.ReadString('\n')
    text = strings.TrimSpace(text)

	for {
		stopChan := make(chan bool)

		go func() {
			_, _ = reader.ReadString('\n')
			stopChan <- true
		}()

		animateText(text, stopChan)

		restartCmd, _ := reader.ReadString('\n')
		if strings.TrimSpace(restartCmd) == "exit" {
			fmt.Println(utils.Color("Goodbye!", utils.Yellow))
			break
		}
	}
}

func animateText(text string, stopChan chan bool) {
	const maxWidth = 40
	pos := 0
	direction := 1
	rainbowColors := []string{utils.Red, utils.Magenta, utils.Yellow, utils.Green, utils.Blue, utils.Black, utils.White}

AnimationLoop:
	for {
		select {
		case <-stopChan:
			break AnimationLoop
		default:
			fmt.Print("\r\033[K")
			coloredText := ""
			for i, char := range text {
				colorIndex := (i + pos) % len(rainbowColors)
				coloredChar := utils.Color(string(char), rainbowColors[colorIndex])
				coloredText += coloredChar
			}

			padding := strings.Repeat(" ", pos)
			fmt.Print(padding + coloredText)

			pos += direction
			if pos == maxWidth-len(text) || pos == 0 {
				direction *= -1
			}

			time.Sleep(100 * time.Millisecond)
		}
	}
}
