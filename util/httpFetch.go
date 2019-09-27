package util

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

//translate util
func HttpFetch(webUrl string) (string, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", webUrl, nil)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var result Result
	if err := json.Unmarshal(body, &result); err != nil {
	}
	return result.Message, nil
}

//获取xpath
func HttpFetchDoc(webUrl string) (*html.Node, error) {

	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://127.0.0.1:8001")
	}
	client := &http.Client{Transport: &http.Transport{Proxy: proxy}}
	req, _ := http.NewRequest("GET", webUrl, nil)
	req.Header.Set("User-Agent", RandomAgent())
	//在解决问题之前需要了解关于go是如何实现connection的一些背景小知识：有两个协程，一个用于读，一个用于写
	// （就是readLoop和writeLoop）。在大多数情况下，readLoop会检测socket是否关闭，并适时关闭connection。
	// 如果一个新请求在readLoop检测到关闭之前就到来了，那么就会产生EOF错误并中断执行，而不是去关闭前一个请求
	// 。这里也是如此，我执行时建立一个新的连接，这段程序执行完成后退出，再次打开执行时服务器并不知道我已经关闭了连接，
	// 所以提示连接被重置；如果我不退出程序而使用for循环多次发送时，旧连接未关闭，新连接却到来，会报EOF。
	req.Close = true
	count := 0
k:
	resp, err := client.Do(req)
	if err != nil {
		if count < 3 {
			fmt.Println("bad request .restart send ")
			count++
			goto k
		} else {
			count = 0
			fmt.Println("request failed: url ---", webUrl)
			return nil, err
		}
	}
	//buff :=new(bytes.Buffer)
	//buff.ReadFrom(resp.Body)
	//s:=buff.String()
	//fmt.Println(s)
	if resp != nil {
		defer resp.Body.Close()
	}

	body, err := htmlquery.Parse(resp.Body)
	if err != nil {
		fmt.Println("解析错误--")
		return nil, err
	}
	return body, nil
}

//random get User-agents
func RandomAgent() string {
	var userAgents = make([]string, 0)
	userAgents = append(userAgents,
		"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (Windows NT 10.0; WOW64; rv:38.0) Gecko/20100101 Firefox/38.0",
		"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; .NET4.0C; .NET4.0E; .NET CLR 2.0.50727; .NET CLR 3.0.30729; .NET CLR 3.5.30729; InfoPath.3; rv:11.0) like Gecko",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)",
		"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0)",
		"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.6; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; TencentTraveler 4.0)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; The World)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; SE 2.X MetaSr 1.0; SE 2.X MetaSr 1.0; .NET CLR 2.0.50727; SE 2.X MetaSr 1.0)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Avant Browser)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1)",
		"Mozilla/5.0 (iPad; U; CPU OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
		"Mozilla/5.0 (Linux; U; Android 2.3.7; en-us; Nexus One Build/FRF91) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
		"MQQBrowser/26 Mozilla/5.0 (Linux; U; Android 2.3.7; zh-cn; MB200 Build/GRJ22; CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
	)
	rand.Seed(time.Now().UnixNano())
	agent := userAgents[rand.Intn(len(userAgents))]

	return agent
}

func ReplaceNumber(ss string) int {

	ko := strings.Replace(ss, "\n", "", -1)
	ls := strings.Replace(ko, "Episode ", "", -1)

	kk := strings.Replace(ls, "\"", "", -1)
	pp := strings.Replace(kk, " ", "", -1)
	number, _ := strconv.Atoi(pp)

	return number
}
