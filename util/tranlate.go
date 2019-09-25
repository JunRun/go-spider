package util

import (
	"fmt"
	"github.com/dollarkillerx/easyutils"
)

func Translate(La string, toLa string, info string) string {
	targetUrl := fmt.Sprintf("http://gsfy.dollarkiller.com/translate?sl=%s&tl=%s&text=%s", La, toLa, info)
	bUrl, err := easyutils.UrlEncoding(targetUrl)
	if err != nil {
		panic("转码失败")
	}
ki:
	bytes, err := HttpFetch(bUrl)
	if err != nil {
		goto ki
	}
	return bytes
}
