package main

import (
	"flag"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nlopes/slack"
)

var config = struct {
	sAPIKey string
	sChat   string
	tAPIKey string
	tChat   int64
	api     *slack.Client
	rtm     *slack.RTM
}{
	sAPIKey: "PLEASE_CHANGE_ME",
	sChat:   "PLEASE_CHANGE_ME",
	tAPIKey: "PLEASE_CHANGE_ME",
	tChat:   0, // PLEASE CHANGE ME
}

func init() {
	sAPIKeyF := flag.String("slack-api-key", "", "SlackBot API KEY")
	sChatF := flag.String("slack-channel", "", "SlackBot Channel ID")
	tAPIKeyF := flag.String("telegram-api-key", "", "Telegram API KEY")
	tChatF := flag.Int64("telegram-channel", 0, "Telegram Channel ID")
	flag.Parse()

	config.sAPIKey = *sAPIKeyF
	config.sChat = *sChatF
	config.tAPIKey = *tAPIKeyF
	config.tChat = *tChatF
}

func main() {
	// Slack client initialize
	config.api = slack.New(config.sAPIKey)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	config.api.SetDebug(false)
	config.rtm = config.api.NewRTM()
	go config.rtm.ManageConnection()

	// Gate from tg to slck
	// Main loop

	// Telegram initialize
	bot, err := tgbotapi.NewBotAPI(config.tAPIKey)
	if err != nil {
		log.Printf("Bot poller: Something wrong with your key, %s", config.tAPIKey)
		log.Panic(err)
	}

	bot.Debug = false
	log.Printf("Bot poller: Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			firstname := update.Message.From.FirstName
			lastname := update.Message.From.LastName
			text := update.Message.Text
			date := update.Message.Date
			currentChatID := update.Message.Chat.ID
			log.Printf("Current chat: %d, from command line: %d", currentChatID, config.tChat)
			if currentChatID != config.tChat {
				log.Printf("Bot poller: Wrong chat: %d FirstName: %s LastName: %s", currentChatID, firstname, lastname)
				continue
			} else {
				log.Println("Bot Poller: Start updating. Updater come in.")
				slckSender(config.rtm, date, text, firstname, lastname)
			}
		}
	}
}

func slckSender(rtm *slack.RTM, date int, text string, firstname string, lastname string) {
	log.Println("Message from bot received")
	message := "Infra Announce: " + text + " " +
		time.Unix(int64(date), 0).Format("Mon Jan _2 15:04") +
		" from " + firstname + " " + lastname
	rtm.SendMessage(rtm.NewOutgoingMessage(message, config.sChat))
	time.Sleep(time.Second)
}
