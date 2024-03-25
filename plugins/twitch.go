package plugins

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"netpuppy/utils"
	twitch "github.com/gempir/go-twitch-irc/v3"
)

func init() {
	Register("twitch", &Twitch{})
}

type Twitch struct{}

func (r *Twitch) Execute(pluginDataChan chan<- string) {
	fmt.Println("[Twitch Chat] Enter channel usernames separated by commas.")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {

		channelNames := strings.Split(scanner.Text(), ",")
		client := twitch.NewAnonymousClient()

		client.OnPrivateMessage(func(message twitch.PrivateMessage) {
			channelColored := utils.Color(fmt.Sprintf("[%s]", message.Channel), utils.Green)
			userColored := utils.Color(message.User.DisplayName, utils.Yellow)
			chatMessage := fmt.Sprintf("%s %s: %s", channelColored, userColored, message.Message)
			pluginDataChan <- chatMessage
		})

		for _, channelName := range channelNames {
			trimmedChannelName := strings.TrimSpace(channelName)
			client.Join(trimmedChannelName)
			fmt.Printf("[Twitch Chat] Connecting to Twitch chat for channel: %s\n", trimmedChannelName)
		}

		err := client.Connect()
		if err != nil {
			fmt.Println("Failed to connect to Twitch IRC:", err)
			return
		}
	} else {
		fmt.Println("Failed to read channel names.")
	}
}
