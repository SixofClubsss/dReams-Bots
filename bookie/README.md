# Telegram-Bot
[Gnomon](https://github.com/civilware/gnomon) powered Bot for relaying bet contract info to Telegram.

### Configure
Create a Telegram bot with the [BotFather](https://core.telegram.org/bots/tutorial).
Make `config.json` file.
```
{
 "botApi": "https://api.telegram.org/bot",
 "apiKey": "TELEGRAM-BOT-TOKEN-HERE",
 "daemon": "127.0.0.1:10102",
 "chatID": yourGroupChatID#,
 "updateConfigs": {
  "limit": 100,
  "timeout": 0,
  "updateFreq": 300000000
 },
 "webhook": false,
 "logFile": "STDOUT",
 "blockedUsers": null
}
```

### Contracts 
You can set custom contracts by changing the contract const in `main.go`
```
btcUSDT_contract  = "c89c2f514300413fd6922c28591196a7c48b42b07e7f4d7d8d9f7643e253a6ff"
deroUSDT_contract = "eaa62b220fa1c411785f43c0c08ec59c761261cb58a0ccedc5b358e5ed2d2c95"
ect...
```

### Run
Install latest [Go](https://go.dev/doc/install) version.

```
git clone https://github.com/SixofClubsss/dReams-Bots.git
cd dReams-Bots
go mod tidy
cd bookie
go run .
```