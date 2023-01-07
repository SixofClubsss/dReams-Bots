package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/SixofClubsss/dReams/menu"
	"github.com/SixofClubsss/dReams/prediction"
	"github.com/SixofClubsss/dReams/rpc"
	"github.com/SixofClubsss/dReams/table"
	"github.com/docopt/docopt-go"
)

type configs struct {
	BotAPI        string `json:"bot_api"`
	APIKey        string `json:"api_key"`
	Daemon        string `json:"daemon_address"`
	UpdateConfigs struct {
		Limit      int `json:"limit"`
		Timeout    int `json:"timeout"`
		UpdateFreq int `json:"update_freq"`
	} `json:"update_configs"`
	Webhook      bool        `json:"webhook"`
	LogFile      string      `json:"log_file"`
	BlockedUsers interface{} `json:"blocked_users"`
}

var command_line string = `dReams Telegram Bot
For relaying bet contract info to Telegram.

Usage:
  Bot [options]
  Bot -h | --help

Options:
  -h --help     Show this screen.
  --debug=<false>	Bot option, true/false value for terminal debug.
  --fastsync=<false>	Gnomon option,  true/false value to define loading at chain height on start up.
  --num-parallel-blocks=<5>   Gnomon option,  defines the number of parallel blocks to index.`

func flags() (debug bool) {
	arguments, err := docopt.ParseArgs(command_line, nil, "v0.9.2")

	if err != nil {
		log.Fatalf("Error while parsing arguments: %s\n", err)
	}

	if arguments["--debug"] != nil {
		if arguments["--debug"].(string) == "true" {
			debug = true
		}
	}

	fastsync := true
	if arguments["--fastsync"] != nil {
		if arguments["--fastsync"].(string) == "false" {
			fastsync = false
		}
	}

	parallel := 1
	if arguments["--num-parallel-blocks"] != nil {
		s := arguments["--num-parallel-blocks"].(string)
		switch s {
		case "2":
			parallel = 2
		case "3":
			parallel = 3
		case "4":
			parallel = 4
		case "5":
			parallel = 5
		default:
			parallel = 1
		}
	}

	menu.Gnomes.Fast = fastsync
	menu.Gnomes.Para = parallel

	return
}

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("")
		StopGnomon(menu.Gnomes.Init)
		log.Println("[Telegram-Bot] Exiting")
		os.Exit(0)
	}()
}

func stamp() {
	fmt.Println(`♤♡♧♢♧♡♤♡♧♢♧♡♤♡♧♢♧♡♤♡♧♢♧♡♤♡♧♢♧♡♤`)
	fmt.Println(`     dReams Telegram Bot`)
	fmt.Println(`   https://dreamtables.net`)
	fmt.Println(`   ©2022-2023 dReam Tables`)
	fmt.Println(`♤♡♧♢♧♡♤♡♧♢♧♡♤♡♧♢♧♡♤♡♧♢♧♡♤♡♧♢♧♡♤`)
}

func readConfig() string {
	if !table.FileExists("configs.json") {
		log.Panicln("[Telegram-Bot] No configs file found")
		return ""
	}

	file, err := os.ReadFile("configs.json")

	if err != nil {
		log.Panicln("[Telegram-Bot]", err)
		return ""
	}

	var config configs
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Panicln("[Telegram-Bot]", err)
		return ""
	}

	rpc.Round.Daemon = config.Daemon

	return config.APIKey
}

func checkConnection() {
	ticker := time.NewTicker(6 * time.Second)
	for range ticker.C {
		rpc.Ping()
		if !rpc.Signal.Daemon {
			log.Println("[Telegram-Bot] Daemon Disconnected")
			StopGnomon(menu.Gnomes.Init)
			log.Println("[Telegram-Bot] Exiting")
			os.Exit(0)
		}
	}
}

func formatFinals(finals []string) (text string) {
	for i := range finals {
		split := strings.Split(finals[i], "   ")
		game := strings.Split(split[1], "_")
		var str string
		if len(game) == 5 {
			str = "Game #" + split[0] + "\n<b>" + game[2] + "</b>" + " Winner: " + prediction.WinningTeam(game[2], game[4])
		} else if len(game) == 4 {
			/// condition until results catch up to v0.9.2 format
			str = "Game #" + split[0] + "\n<b>" + game[1] + "</b>" + " Winner: " + prediction.WinningTeam(game[1], game[3])
		} else {
			str = "Game #" + split[0] + "\n<b>" + game[2] + "</b>" + " Tie"
		}
		text = text + "\n\n" + str + "\nTXID: <code>" + split[2] + "</code>"
	}

	return
}
func HtmlTextParse(text string) (new string) { /// format bot reply for html
	first := parseForID(text)
	second := parseForTime(first)
	new = parseForPrediction(second)

	return
}

func rangeCheck(s string, r []string) bool {
	for i := range r {
		if s == r[i] {
			return false
		}
	}

	return true
}

func parseForID(text string) (new string) { /// format scids and txids as code
	var more bool
	var found []string
	sc, _ := regexp.Compile(`^\w{64,64}$`)

	for i := range text {
		if i == len(text)-63 {
			break
		}

		if sc.MatchString(text[i : i+64]) {
			found = append(found, text[i:i+64])
			new = text[:i] + "<code>" + text[i:i+64] + "</code>" + text[i+64:]
			break
		}
	}

	if new == "" {
		new = text
		return
	}

	for {
		for i := range new {
			if i == len(new)-63 {
				break
			}

			if sc.MatchString(new[i:i+64]) && rangeCheck(new[i:i+64], found) {
				more = true
				found = append(found, new[i:i+64])
				new = new[:i] + "<code>" + new[i:i+64] + "</code>" + new[i+64:]
				break
			}
			more = false
		}

		if !more {
			break
		}

	}

	return
}

func parseForPrediction(text string) (new string) { /// format predictions bold
	pre, _ := regexp.Compile(`^\w{3,4}\-\w{3,4}$|(DERO-Block Time)|(DERO-Block Number)|(DERO-Difficulty)`)

	for i := range text {
		if i == len(text)-7 {
			break
		}

		if pre.MatchString(text[i : i+8]) {
			if text[i:i+15] == "DERO-Block Time" || text[i:i+15] == "DERO-Difficulty" {
				new = text[:i] + "<b>" + text[i:i+15] + "</b>" + text[i+15:]
				break
			} else if text[i:i+17] == "DERO-Block Number" {
				new = text[:i] + "<b>" + text[i:i+17] + "</b>" + text[i+17:]
				break
			} else {
				if text[i:i+4] == "DERO" {
					new = text[:i] + "<b>" + text[i:i+9] + "</b>" + text[i+9:]
					break
				}
				new = text[:i] + "<b>" + text[i:i+8] + "</b>" + text[i+8:]
				break
			}
		}
	}

	if new == "" {
		new = text
		return
	}

	return
}

func parseForTime(text string) (new string) { /// format time to utc
	var more bool
	t, _ := regexp.Compile(`^\d{4,4}\-\d{2,2}\-\d{2,2}\ \d{2,2}\:\d{2,2}\:\d{2,2}\ \-\d{4,4}\ \w{3,3}$`)

	for i := range text {
		if i == len(text)-28 {
			break
		}

		if t.MatchString(text[i : i+29]) {
			rc3339 := text[i:i+10] + "T" + text[i+11:i+19] + text[i+20:i+23] + ":" + text[i+23:i+25]
			format := time.RFC3339
			utc, err := time.Parse(format, rc3339)
			if err != nil {
				fmt.Println(err)
			}
			new = text[:i] + utc.UTC().String() + text[i+29:]
			break
		}
	}

	if new == "" {
		new = text
		return
	}

	for {
		for i := range new {
			if i == len(new)-28 {
				break
			}

			if t.MatchString(new[i : i+29]) {
				rc3339 := new[i:i+10] + "T" + new[i+11:i+19] + new[i+20:i+23] + ":" + new[i+23:i+25]
				format := time.RFC3339
				utc, err := time.Parse(format, rc3339)
				if err != nil {
					fmt.Println(err)
				}
				more = true
				new = new[:i] + utc.UTC().String() + new[i+29:]
				break
			}
			more = false
		}

		if !more {
			break
		}
	}

	return
}
