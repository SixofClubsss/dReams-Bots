package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/SixofClubsss/dReams/menu"
	"github.com/SixofClubsss/dReams/rpc"
	"github.com/civilware/Gnomon/indexer"
)

func searchFilters() (filter []string) {
	predict, _ := rpc.GetPredictCode(rpc.Signal.Daemon, 0)
	if predict != "" {
		filter = append(filter, predict)
	}

	sports, _ := rpc.GetSportsCode(rpc.Signal.Daemon, 0)
	if sports != "" {
		filter = append(filter, sports)
	}

	return
}

func startGnomon() {
	log.Println("[Telegram-Bot] Starting Gnomon")
	backend := menu.GnomonDB()

	last_height := backend.GetLastIndexHeight()
	daemon_endpoint := rpc.Round.Daemon
	runmode := "daemon"
	mbl := false
	closeondisconnect := false
	filter := searchFilters()

	if len(filter) == 2 {
		menu.Gnomes.Indexer = indexer.NewIndexer(backend, filter, last_height, daemon_endpoint, runmode, mbl, closeondisconnect, menu.Gnomes.Fast)
		go menu.Gnomes.Indexer.StartDaemonMode(menu.Gnomes.Para)
		time.Sleep(3 * time.Second)
		menu.Gnomes.Init = true
		menu.Gnomes.Sync = true
	}
}

func StopGnomon(gi bool) {
	if gi && !menu.GnomonClosing() {
		log.Println("[Telegram-Bot] Putting Gnomon to Sleep")
		menu.Gnomes.Indexer.Close()
		menu.Gnomes.Init = false
		time.Sleep(1 * time.Second)
		log.Println("[Telegram-Bot] Gnomon is Sleeping")
	}
}

func GetBook(dc bool, scid string) (text string) {
	if dc && !menu.GnomonClosing() && !menu.GnomonWriting() {
		_, initValue := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "s_init", menu.Gnomes.Indexer.ChainHeight, true)
		if initValue != nil {
			init := initValue[0]
			var single bool
			iv := 1
			for {
				_, s_init := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "s_init_"+strconv.Itoa(iv), menu.Gnomes.Indexer.ChainHeight, true)
				if s_init != nil {
					game, _ := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "game_"+strconv.Itoa(iv), menu.Gnomes.Indexer.ChainHeight, true)
					league, _ := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "league_"+strconv.Itoa(iv), menu.Gnomes.Indexer.ChainHeight, true)
					//_, s_n := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "s_#_", menu.Gnomes.Indexer.ChainHeight, true)
					_, s_amt := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "s_amount_"+strconv.Itoa(iv), menu.Gnomes.Indexer.ChainHeight, true)
					_, s_end := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "s_end_at_"+strconv.Itoa(iv), menu.Gnomes.Indexer.ChainHeight, true)
					_, s_total := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "s_total_"+strconv.Itoa(iv), menu.Gnomes.Indexer.ChainHeight, true)
					//s_urlValue, _ := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "s_url_", menu.Gnomes.Indexer.ChainHeight, true)
					_, s_ta := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "team_a_"+strconv.Itoa(iv), menu.Gnomes.Indexer.ChainHeight, true)
					_, s_tb := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "team_b_"+strconv.Itoa(iv), menu.Gnomes.Indexer.ChainHeight, true)
					//_, time_a := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "time_a", menu.Gnomes.Indexer.ChainHeight, true)
					// _, time_b := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "time_b", menu.Gnomes.Indexer.ChainHeight, true)

					team_a := menu.TrimTeamA(game[0])
					team_b := menu.TrimTeamB(game[0])

					eA := time.Unix(int64(s_end[0]), 0).UTC()
					closes := fmt.Sprint(eA.Format("2006-01-02 15:04:05 UTC"))

					min := fmt.Sprint(float64(s_amt[0]) / 100000)

					aV := strconv.Itoa(int(s_ta[0]))
					bV := strconv.Itoa(int(s_tb[0]))

					float := float64(s_total[0])
					total := fmt.Sprint(float / 100000)

					now := time.Now().UTC()
					nowTime := fmt.Sprint(now.Format("2006-01-02 15:04:05 UTC"))
					if !single {
						single = true
						text = "SCID: \n<code>" + scid + "</code>\n\n" + "Time now is: " + nowTime + "\n\nCurrent Games:\n\n"
					}

					var pre string
					if league[0] == "UFC" || league[0] == "Bellator" {
						pre = "Fight #"
					} else {
						pre = "Game #"
					}

					text = text + pre + strconv.Itoa(iv) + "\nLeague: " + league[0] + "\n<b>" + game[0] + "</b>\nMinimum: " + min + " Dero\nPot: " + total + "\nCloses at: " + closes + "\n" + team_a + " Picks: " + aV + "\n" + team_b + " Picks: " + bV + "\n\n"

				}

				if iv >= int(init) {
					break
				}

				iv++
			}

			if text == "" {
				text = "No Results Found"
			}
			return
		}
	}

	return
}
