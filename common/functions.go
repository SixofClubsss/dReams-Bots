package common

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/civilware/Gnomon/structures"
	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/dReam-dApps/dReams/rpc"
	"github.com/docopt/docopt-go"
	"github.com/sirupsen/logrus"
)

type config struct {
	BotAPI        string `json:"botApi"`
	APIKey        string `json:"apiKey"`
	Daemon        string `json:"daemon"`
	SCID          string `json:"scid"`
	ChatID        int64  `json:"chatID"`
	UpdateConfigs struct {
		Limit      int `json:"limit"`
		Timeout    int `json:"timeout"`
		UpdateFreq int `json:"updateFreq"`
	} `json:"updateConfigs"`
	Webhook      bool        `json:"webhook"`
	LogFile      string      `json:"logFile"`
	BlockedUsers interface{} `json:"blockedUsers"`
}

var gnomon = gnomes.NewGnomes()
var logger = structures.Logger.WithFields(logrus.Fields{})

var commands string = `dReams Telegram Bot
For relaying blockchain data to Telegram.

Usage:
  Bot [options]
  Bot -h | --help

Options:
  -h --help     Show this screen.
  --debug=<false>	Bot option, true/false value for terminal debug.
  --fastsync=<false>	Gnomon option,  true/false value to define loading at chain height on start up.
  --num-parallel-blocks=<5>   Gnomon option,  defines the number of parallel blocks to index.`

func Flags() (debug bool) {
	arguments, err := docopt.ParseArgs(commands, nil, rpc.Version().String())
	if err != nil {
		logger.Fatalf("Error while parsing arguments: %s\n", err)
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

	gnomon.SetDBStorageType("boltdb")
	gnomon.SetFastsync(fastsync, true, 10000)
	gnomon.SetParallel(parallel)

	fmt.Println(`♤♡♧♢♧♡♤♡♧♢♧♡♤♡♧♢♧♡♤♡♧♢♧♡♤♡♧♢♧♡♤`)
	fmt.Println(`      dReams Telegram Bot`)
	fmt.Println(`     https://dreamdapps.io`)
	fmt.Println(`       ©2023 SixofClubs`)
	fmt.Println(`♤♡♧♢♧♡♤♡♧♢♧♡♤♡♧♢♧♡♤♡♧♢♧♡♤♡♧♢♧♡♤`)

	return
}

func ReadConfig(tag string) (c config) {
	if !dreams.FileExists("config.json", tag) {
		logger.Fatalf("[%s] No config file found\n", tag)
	}

	file, err := os.ReadFile("config.json")
	if err != nil {
		logger.Fatalf("[%s] %s\n", tag, err)
	}

	err = json.Unmarshal(file, &c)
	if err != nil {
		logger.Fatalf("[%s] %s\n", tag, err)
	}

	rpc.Daemon.Rpc = c.Daemon

	return
}

var done = make(chan struct{})

// Kill the connection loop
func Close() {
	close(done)
}

// Check daemon connection
func CheckConnection(tag string) {
	ticker := time.NewTicker(6 * time.Second)
	for {
		select {
		case <-ticker.C:
			rpc.Ping()
			if !rpc.Daemon.IsConnected() {
				logger.Printf("[%s] Daemon Disconnected\n", tag)
				gnomon.Stop(tag)
				ticker.Stop()
				logger.Printf("[%s] Exiting\n", tag)
				os.Exit(0)
			} else {
				gnomes.State(false, nil)
			}

		case <-done:
			ticker.Stop()
			return
		}
	}
}
