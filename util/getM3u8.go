package util

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

type M3u8 struct {
	hardsubLang string `json:"hardsub_lang"`
	url         string `json:"url"`
}

//正则匹配 寻找.m3u8链接
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

func OpenFile() {
	file, err := os.OpenFile("./log.txt", os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(file)
	write := bufio.NewWriter(file)
	write.WriteString("lllll\n")
	write.Flush()
	file.Close()
}
