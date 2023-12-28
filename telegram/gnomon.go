package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/SixofClubsss/dPrediction/prediction"
	"github.com/civilware/Gnomon/structures"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/dReam-dApps/dReams/rpc"
	"github.com/sirupsen/logrus"
)

var gnomon = gnomes.NewGnomes()
var logger = structures.Logger.WithFields(logrus.Fields{})

func searchFilters() (filter []string) {
	predict := prediction.GetPredictCode(0)
	if predict != "" {
		filter = append(filter, predict)
	}

	sports := prediction.GetSportsCode(0)
	if sports != "" {
		filter = append(filter, sports)
	}

	return
}

func GetBook(scid string) (text string) {
	if rpc.Daemon.IsConnected() && !gnomon.IsClosing() && !gnomon.IsWriting() {
		_, initValue := gnomon.GetSCIDValuesByKey(scid, "s_init")
		if initValue != nil {
			init := initValue[0]
			var single bool
			iv := 1
			for {
				_, s_init := gnomon.GetSCIDValuesByKey(scid, "s_init_"+strconv.Itoa(iv))
				if s_init != nil {
					game, _ := gnomon.GetSCIDValuesByKey(scid, "game_"+strconv.Itoa(iv))
					league, _ := gnomon.GetSCIDValuesByKey(scid, "league_"+strconv.Itoa(iv))
					//_, s_n := gnomon.GetSCIDValuesByKey(scid, "s_#_")
					_, s_amt := gnomon.GetSCIDValuesByKey(scid, "s_amount_"+strconv.Itoa(iv))
					_, s_end := gnomon.GetSCIDValuesByKey(scid, "s_end_at_"+strconv.Itoa(iv))
					_, s_total := gnomon.GetSCIDValuesByKey(scid, "s_total_"+strconv.Itoa(iv))
					//s_urlValue, _ := gnomon.GetSCIDValuesByKey(scid, "s_url_")
					_, s_ta := gnomon.GetSCIDValuesByKey(scid, "team_a_"+strconv.Itoa(iv))
					_, s_tb := gnomon.GetSCIDValuesByKey(scid, "team_b_"+strconv.Itoa(iv))
					//_, time_a := gnomon.GetSCIDValuesByKey(scid, "time_a")
					// _, time_b := gnomon.GetSCIDValuesByKey(scid, "time_b")

					team_a := prediction.TrimTeamA(game[0])
					team_b := prediction.TrimTeamB(game[0])

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
