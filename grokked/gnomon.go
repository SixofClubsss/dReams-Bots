package main

import (
	"fmt"
	"time"

	"github.com/civilware/Gnomon/structures"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/sirupsen/logrus"
)

var gnomon = gnomes.NewGnomes()
var logger = structures.Logger.WithFields(logrus.Fields{})

// Get the current Grok on scid and return reply message text
func GetGrok(scid string) string {
	if _, grok := gnomon.GetSCIDValuesByKey(scid, "grok"); grok != nil {
		if addr, _ := gnomon.GetSCIDValuesByKey(scid, grok[0]); addr != nil {
			left := "I am not sure how much time is left? contact the dev"
			now := uint64(time.Now().Unix())
			_, last := gnomon.GetSCIDValuesByKey(scid, "last")
			_, dur := gnomon.GetSCIDValuesByKey(scid, "duration")
			if last != nil && dur != nil {
				tf := last[0] + dur[0]
				if now < tf {
					left = fmt.Sprintf("<b>%d</b> minutes left to pass", (tf-now)/60)
				} else if tf != 0 {
					left = fmt.Sprintf("<b>%d</b> minutes past", (now-tf)/60)
				}
			}

			return fmt.Sprintf("Grok is currently:\n\n<code>%s</code>\n\n%s", addr[0], left)
		}
	}

	return "I am not sure? contact the dev"
}
