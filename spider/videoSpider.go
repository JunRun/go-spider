package spider

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/rs/xid"
	"go-spider/modle"
	"go-spider/util"
	"golang.org/x/net/html"
	"strconv"
	"strings"
	"sync"
)

var (
	webUrl   = "https://www.crunchyroll.com/videos/anime/popular/ajax_page?pg="
	videoUrl = "https://www.crunchyroll.com"
	pool     = make(chan int, 20)
	sy1      = sync.WaitGroup{}
	lock     = sync.Mutex{}
	count    = 1
)

func VideoUrlSpider(url string) {
	html, err := util.HttpFetchDoc(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	//获取连接
	nodes := htmlquery.Find(html, "//li/div/a")
	for _, node := range nodes {
		if node == nil {
			break
		}

		videoPath := htmlquery.FindOne(node, "./@href").FirstChild.Data
		videoName := htmlquery.FindOne(node, "./@title").FirstChild.Data
		video := VideoInfoSpider(videoUrl + videoPath)

		video.VideoName = videoName
		o := modle.GetConnect()
		//插入信息
		e := o.Create(&video).Error
		if e != nil {
			fmt.Println("动漫信息 插入失败\n", videoName)
		}
		fmt.Println("动漫信息 插入成功\n", videoName)
	}
}

// 获取动漫信息
func VideoInfoSpider(url string) modle.Video {
	htmlContent, err := util.HttpFetchDoc(url)
	video := modle.Video{}
	if err != nil {
		fmt.Println(err)
		return video
	}

	//获取动漫视频 图片信息
	if image := htmlquery.FindOne(htmlContent, "//*[@id='sidebar_elements']/li[1]/img/@src"); image != nil {
		video.Img = image.FirstChild.Data
	}

	//获取动漫视频 描述
	if info := htmlquery.FindOne(htmlContent, "//*[@id='sidebar_elements']/li/p/span[@class='more']"); info != nil {
		target := info.FirstChild.Data
		target = strings.TrimPrefix(target, "\n")
		target = strings.TrimSuffix(target, "\n")
		target = strings.TrimSpace(target)
		video.Describe = target
	}

	//获取动漫视频 评分
	if mark := htmlquery.FindOne(htmlContent, "//*[@id='showview_about_rate_widget']/@content"); mark != nil {
		video.Mark, _ = strconv.ParseFloat(mark.FirstChild.Data, 64)
	}

	//获取动漫视频 标签
	if tags := htmlquery.Find(htmlContent, "//*[@id='sidebar_elements']/li/ul/li/a"); tags != nil {
		array := new(strings.Builder)
		for _, tag := range tags {
			tagInfo := tag.FirstChild.Data
			array.WriteString(tagInfo + ",")
		}
		video.Tag = array.String()
	}

	video.Vid = xid.New().String()
	episodeList := htmlquery.Find(htmlContent, "//*[@id='showview_content_videos']/ul/li/ul/li/div/a")
	for _, episode := range episodeList {
		if episode == nil {
			break
		}
		pool <- 1
		go func(es *html.Node) {
			defer func() {
				<-pool
			}()
			episodeUrl := htmlquery.FindOne(es, "./@href").FirstChild.Data
			title := htmlquery.FindOne(es, "./img/@alt").FirstChild.Data
			img := ""
			if imgNode := htmlquery.FindOne(es, "./img/@src"); imgNode != nil {
				img = imgNode.FirstChild.Data
			} else if imgNode2 := htmlquery.FindOne(es, "./img/@data-thumbnailurl"); imgNode2 != nil {
				img = imgNode2.FirstChild.Data
			} else {
				img = "null-wide.jpg"
			}
			episodeNumber := htmlquery.FindOne(es, "./span/text()").Data
			EpisodeInfo(videoUrl+episodeUrl, title, img, episodeNumber, video.Vid)
			return
		}(episode)

	}
	//episodeNodes :=htmlquery.Find(html,"")

	return video
}

func EpisodeInfo(url string, title string, img string, episodeNumber string, videoId string) {
	html, err := util.HttpFetchDoc(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	episode := modle.Episode{}
	episode.Vid = videoId
	episode.EpisodeNumber = util.ReplaceNumber(episodeNumber)
	episode.EpisodeName = title
	episode.Img = img
	episode.PlayImg = img[:len(img)-8] + "fwide.jpg"
	episode.UrlPath = url

	if season := htmlquery.FindOne(html, "//*[@id='template_body']/div[3]/div[1]/div/h1/a/span"); season != nil {
		episode.Season = season.FirstChild.Data
	}

	if describe := htmlquery.FindOne(html, "//*[@id='showmedia_about_info']/p/span[1]"); describe != nil {
		episode.Describe = describe.FirstChild.Data
	}

	if languageList := htmlquery.Find(html, "//*[@id='showmedia_about_info_details']/div[1]/span/img/@src"); languageList != nil {
		sb := strings.Builder{}
		for _, languageNode := range languageList {
			country := languageNode.FirstChild.Data
			s := country[len(country)-6 : len(country)-4]
			sb.WriteString(s + "|")
		}
		episode.Language = sb.String()
	}

	if content := htmlquery.FindOne(html, "//*[@id='showmedia_video_box']/script[3]/text()"); content != nil {
		episode.Content = util.GetM3u8(content.Data)
	}

	if tagNode := htmlquery.Find(html, "//*[@id='showmedia_about_info_details']/div[5]/span/a/text()"); tagNode != nil {
		sb := strings.Builder{}
		for _, tag := range tagNode {
			sb.WriteString(tag.Data)
		}
		episode.Tag = sb.String()
	}
	//插入剧集信息
	g := modle.GetConnect()
	if g.Create(&episode).Error != nil {
		fmt.Println("插入失败==== ", episode.EpisodeName)
	}
	fmt.Println("插入成功=== ", episode.EpisodeName)
	return
}

func Run() {
	page := 0
	for {
		if page > 30 {
			break
		}
		VideoUrlSpider(webUrl + strconv.Itoa(page))
		page = page + 1
		fmt.Println("===", page)
	}

}

func UpdateM3u8Task() {
	g := modle.GetConnect()
	var episodes []modle.Episode
	g.Find(&episodes)
	fmt.Println(len(episodes))
	updatePool := make(chan int, 20)
	for _, node := range episodes {
		updatePool <- 1
		go func(episode modle.Episode) {
			defer func() {
				<-updatePool
			}()
			UpdateM3u8(episode)
		}(node)
	}
	close(updatePool)

}

func UpdateM3u8(episode modle.Episode) {
	html, err := util.HttpFetchDoc(episode.UrlPath)
	if err != nil {
		fmt.Println(err)
	}
	if content := htmlquery.FindOne(html, "//*[@id='showmedia_video_box']/script[3]/text()"); content != nil {
		episode.Content = util.GetM3u8(content.Data)
	}
	db := modle.GetConnect()
	db.Save(&episode)
	fmt.Println("更新成功----", episode.EpisodeName)
}
