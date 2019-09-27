package test

import (
	"encoding/json"
	"fmt"
	"github.com/arl/assertgo"
	"github.com/dollarkillerx/easyutils"
	"github.com/robfig/cron/v3"
	"go-spider/modle"
	"go-spider/spider"
	"go-spider/util"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestP(t *testing.T) {

	targetUrl := fmt.Sprintf("https://gsfy.dollarkiller.com/translate?sl=%s&tl=%s&text=%s", "", "en", "你好")
	bUrl, err := easyutils.UrlEncoding(targetUrl)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", bUrl, nil)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	var result util.Result
	if err := json.Unmarshal(body, &result); err != nil {

	}
	fmt.Println("===============" + result.Message)
	//bytes, err := httplib.EuUserGet(bUrl)
	//if err != nil {
	//	panic(err)
	//}
	//
	//clog.Println(string(bytes))

}

func TestSingle(t *testing.T) {
	db := modle.GetConnect()
	//var i int
	fmt.Printf("db:%v\n", db)
	ls := modle.GetConnect()
	fmt.Printf("ls:%v\n", ls)
	fmt.Println(ls == db)
}

func TestSpider(t *testing.T) {
	var webUr = "https://www.crunchyroll.com/videos/anime/popular/ajax_page?pg=4"
	spider.VideoUrlSpider(webUr)
}

func TestVideo(t *testing.T) {
	var wurl = "https://www.crunchyroll.com/naruto-shippuden"
	spider.VideoInfoSpider(wurl)
}

func TestV(t *testing.T) {
	var m = modle.Video{}
	assert.False(false, m)

}

func TestGetM3u8(t *testing.T) {
	spider.EpisodeInfo("https://www.crunchyroll.com/naruto-shippuden/episode-495-hidden-leaf-story-the-perfect-day-for-a-wedding-part-2-a-full-powered-wedding-gift-727637", "test", "test", "32", "31")
}

func TestStringSlice(t *testing.T) {
	//s:="https://img1.ak.crunchyroll.com/i/spire2-tmb/b37856fbefcc27c27f87794164534ecc1490206055_wide.jpg"
	//imp :=s[:len(s)-8]+"fwide.jpg"
	ss := `"
	Episode 480
	"`
	ko := strings.Replace(ss, "\n\t", "", -1)
	ls := strings.Replace(ko, "Episode ", "", -1)
	kk := strings.Replace(ls, "\"", "", -1)
	number, err := strconv.Atoi(kk)
	if err != nil {
		panic(err)
	}
	fmt.Println(number)
}

func TestUpdateM3u8(t *testing.T) {
	spider.UpdateM3u8Task()
}

func TestCorn(t *testing.T) {
	fmt.Println("定时任务开启")
	c := cron.New()
	_, err := c.AddFunc("*/1 * * * *", spider.UpdateM3u8Task)
	if err != nil {
		panic(err)
	}
	c.Start()
	select {}
}

func TestFile(t *testing.T) {
	util.OpenFile()
}

var pool = make(chan int, 100)
var ch = make(chan struct{})
var sy = sync.WaitGroup{}

func TestLanguage(t *testing.T) {
	var videoList []modle.Video
	count := 1
	for {
		db := modle.GetConnect()
		db.Limit(1).Offset((count - 1) * 1).Find(&videoList)
		//UpdateLanguage(videoList)
		pool <- 1
		sy.Add(1)
		go func(list []modle.Video) {
			defer func() {
				sy.Done()
				<-pool
			}()
			UpdateLanguage(list)

		}(videoList)
		count++
		if count >= 1040 {
			break
		}

	}
	sy.Wait()
	//time.Sleep(time.Second*10)
}

func UpdateLanguage(videoList []modle.Video) {
	db := modle.GetConnect()
	for _, video := range videoList {
		//pool<-1
		//sy.Add(1)
		//go func(video modle.Video) {
		//	defer func() {
		//		sy.Done()
		//		<-pool
		//	}()
		var episode modle.Episode
		if err := db.Where("vid = ?", video.Vid).First(&episode).Error; err != nil {
			fmt.Println("查询失败", video.Vid, err)
		}
		video.Language = episode.Language

		if db.Save(video).Error != nil {
			fmt.Println("保存失败", video.Vid)
		}
		//}(video)

	}
	//if len(videoList)<100{
	//	close(ch)
	//}
}

func TestPP(t *testing.T) {

}
