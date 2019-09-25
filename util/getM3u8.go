package util

import (
	"regexp"
	"strings"
)

type M3u8 struct {
	hardsubLang string `json:"hardsub_lang"`
	url         string `json:"url"`
}

func GetM3u8(text string) string {
	r := regexp.MustCompile(`\{"format":([\w\W]*?)\.m3u8([\w\W]*?)\}`)
	machs := r.FindAllString(text, -1)

	sb := strings.Builder{}
	sb.WriteString("[")
	for _, s := range machs {
		sb.WriteString(s)
	}
	sb.WriteString("]")
	return sb.String()
}
