package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"
)

var config = struct {
	sAPI  string
	sChat string
	tAPI  string
	tChat string
}{
	sAPI:  "PLEASE_CHANGE_ME",
	sChat: "PLEASE_CHANGE_ME",
	tAPI:  "PLEASE_CHANGE_ME",
	tChat: "PLEASE_CHANGE_ME",
}

func init() {
	sAPIF := flag.String("slack-api-key", "", "SlackBot API KEY")
	sChatF := flag.String("slack-channel", "", "SlackBot Channel ID")
	tAPIF := flag.String("telegram-api-key", "", "Telegram API KEY")
	tChatF := flag.String("telegram-channel", "", "Telegram Channel ID")
	flag.Parse()

	config.sAPI = *sAPIF
	config.sChat = *sChatF
	config.tAPI = *tAPIF
	config.tChat = *tChatF
}

func main() {
	api := slack.New(config.sAPI)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()

	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace C94CCHD8A with your Channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", config.sChat))

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)
			rtm.SendMessage(rtm.NewOutgoingMessage("Test", config.sChat))

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}
