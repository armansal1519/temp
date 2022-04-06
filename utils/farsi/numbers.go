package farsi

import (
	"strings"
)

func Translate(initial string) (string, string) {
	t := make(map[string]string)
	t["1"] = "۱"
	t["2"] = "۲"
	t["3"] = "۳"
	t["4"] = "۴"
	t["5"] = "۵"
	t["6"] = "۶"
	t["7"] = "۷"
	t["8"] = "۸"
	t["9"] = "۹"
	t["0"] = "۰"

	farsiStr := strings.Clone(initial)
	englishStr := strings.Clone(initial)

	for s, s2 := range t {
		farsiStr = strings.ReplaceAll(farsiStr, s, s2)
	}
	for s, s2 := range t {
		englishStr = strings.ReplaceAll(englishStr, s2, s)
	}

	return farsiStr, englishStr
}
