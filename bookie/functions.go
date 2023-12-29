package main

import (
	"regexp"
	"strings"

	"github.com/SixofClubsss/dPrediction/prediction"
	"github.com/SixofClubsss/dReams-Bots/common"
)

// Format sports reply message text
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

// Format prediction reply message text
func formatPrediction(text string) (new string) {
	result := common.HtmlTextParse(text)

	pre, _ := regexp.Compile(`^\w{3,4}\-\w{3,4}$|(DERO-Block Time)|(DERO-Block Number)|(DERO-Difficulty)`)

	for i := range result {
		if i == len(result)-7 {
			break
		}

		if pre.MatchString(result[i : i+8]) {
			if result[i:i+15] == "DERO-Block Time" || result[i:i+15] == "DERO-Difficulty" {
				new = result[:i] + "<b>" + result[i:i+15] + "</b>" + result[i+15:]
				break
			} else if result[i:i+17] == "DERO-Block Number" {
				new = result[:i] + "<b>" + result[i:i+17] + "</b>" + result[i+17:]
				break
			} else {
				if result[i:i+4] == "DERO" {
					new = result[:i] + "<b>" + result[i:i+9] + "</b>" + result[i+9:]
					break
				}
				new = result[:i] + "<b>" + result[i:i+8] + "</b>" + result[i+8:]
				break
			}
		}
	}

	if new == "" {
		new = result
		return
	}

	return
}
