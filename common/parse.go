package common

import (
	"regexp"
	"time"
)

// Parse text for HTML tags
func HtmlTextParse(text string) string {
	return parseForTime(parseForSCID(text))
}

// Check that we have not already tagged this string
func haveAlready(s string, r []string) bool {
	for i := range r {
		if s == r[i] {
			return true
		}
	}

	return false
}

// Parse for scid(s) and txid(s) and add HTML code tags if found
func parseForSCID(text string) (new string) {
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

			if sc.MatchString(new[i:i+64]) && !haveAlready(new[i:i+64], found) {
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

// Parse for time and format it to RFC3339 and UTC zone
func parseForTime(text string) (new string) {
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
				logger.Errorln("[parseForTime]", err)
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
					logger.Errorln("[parseForTime]", err)
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
