# tg2slck

### Telegram to Slack bridge
Brings messages from Telegram chat to Slack channel. Simple but efficient.

### Usage:
```
https://github.com/doctornkz/tg2slck.git
./tg2slck -slack-api-key=<SLACK_BOT_API_KEY> -slack-channel=<SLACK_CHANNEL_ID>  -telegram-api-key=<TELEGRAM_BOT_API_KEY> -telegram-channel=<TELEGRAM_CHANNEL

2018/02/07 00:25:10 Bot poller: Authorized on account ****Bot
2018/02/07 00:25:10 Bot poller: Pre-update section
2018/02/07 00:25:15 Current chat: ******, from command line: *****
2018/02/07 00:25:15 Bot Poller: Start updating. Updater come in.
2018/02/07 00:25:15 Message from bot received
2018/02/07 00:25:16 Bot poller: Pre-update section
```

### TODO: 
1. Implement goroutines
2. Simplify code
3. Make Dockerfile, Jenkinsfile and Makefile
4. Push Docker image
