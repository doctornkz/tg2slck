package main

import (
	"flag"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/slack-go/slack"
)

// Configuration and defaults for hardcoded binary mode
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

// TODO: Support full mesh syncronizing:  API_1+Chat_1, API_2+Chat_2,[API_N+Chat_N]...

// ChanPost - message structure, some methods are not mirrored from Message, beware.
type ChanPost struct {
	date int
	text string
	/*
		firstname string  // In Bot API 3.5 From  methods //
		lastname  string  //     ===  UNSUPPORTED ===     //
	*/
}

// Command-line processing
func init() {
	sAPIKeyF := flag.String("slack-api-key", config.sAPIKey, "SlackBot API KEY")
	sChatF := flag.String("slack-channel", config.sChat, "SlackBot Channel ID")
	tAPIKeyF := flag.String("telegram-api-key", config.tAPIKey, "Telegram API KEY")
	tChatF := flag.Int64("telegram-channel", config.tChat, " Telegram Channel ID (default 0)")
	flag.Parse()

	config.sAPIKey = *sAPIKeyF
	config.sChat = *sChatF
	config.tAPIKey = *tAPIKeyF
	config.tChat = *tChatF
}

func main() {
	// Slack client initialize
	config.api = slack.New(config.sAPIKey)
	config.rtm = config.api.NewRTM()
	go config.rtm.ManageConnection()

	/*

		Main loop based on infinite FOR around telegram updates.ChanPost
		 If you want to have duplex sync, change it to go routine logic.

	*/

	// Telegram initialize
	bot, err := tgbotapi.NewBotAPI(config.tAPIKey)
	if err != nil {
		log.Printf("Check Telegram API Key: %v", err)
		os.Exit(1)
	}

	bot.Debug = false
	log.Printf("Bot poller: Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Printf("Updates error: %v", err)
		os.Exit(1)
	}

	/*
		ChannelPost is Message's equivalent for
		channel management. To use Chat-Mode,
		feel free to replace ChannetPost to Message.
	*/

	for {
		select {
		case update := <-updates:
			if update.ChannelPost == nil {
				log.Println("Nil message, continue")
				continue
			}
			// Structure updating
			cp := ChanPost{date: update.ChannelPost.Date, text: update.ChannelPost.Text}
			// Check permission. Bot send nothing to irrelevant chat.
			currentChatID := update.ChannelPost.Chat.ID
			log.Printf("Current chat: %d, from command line: %d", currentChatID, config.tChat)
			if currentChatID != config.tChat {
				log.Printf("Bot poller: Wrong chat: %d ", currentChatID)
				continue
			} else {
				log.Println("Bot Poller: Send message to slack channel")
				cp.slckSend(config.rtm)
			}
		}
	}
}

// Sending to Slack channel function.
func (c *ChanPost) slckSend(rtm *slack.RTM) {
	log.Println("Message from bot received")
	message := "Infra Announce: " + c.text + " " +
		time.Unix(int64(c.date), 0).Format("Mon Jan _2 15:04") +
		"(from telegram channel)"
	rtm.SendMessage(rtm.NewOutgoingMessage(message, config.sChat))
}
