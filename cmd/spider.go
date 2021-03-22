package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/linuxing3/vpsman/util"
	"github.com/spf13/cobra"
)

var (
	DocBodySelector = "document.querySelector('body')"
	webSite         = ""
	keyWord         = ""
)

// spiderCmd represents the spider command
var spiderCmd = &cobra.Command{
	Use:   "spider",
	Short: "Go powered Spider",
	Long: ` About usage of using spider. For example: 
Gospider is a CLI Command for Go that crawl web contents.
to quickly create records in postgressql.`,
	Run: func(cmd *cobra.Command, args []string) {
		if webSite == "douban" {
			scrapDouban()
		} else if webSite == "googlenews" {
			scrapGoogleNews(keyWord)
		} else if webSite == "iciba" {
			scrapIciba()
		} else {
			spiderMenu()
		}
	},
}

// TrojanMenu 控制TrojanMenu
func spiderMenu() {
exit:
	for {
		fmt.Println()
		fmt.Print(util.Cyan("请选择"))
		fmt.Println()
		loopMenu := []string{"iciba", "douban", "hacknews", "googlenews"}
		choice := util.LoopInput("回车退出:   ", loopMenu, false)
		switch choice {
		case 1:
			script := "docker exec -it spider /root/go/bin/gospider spider"
			util.ExecCommand(script)
		case 2:
			script := "docker exec -it spider /root/go/bin/gospider spider --website=douban"
			util.ExecCommand(script)
		case 3:
			script := "docker exec -it spider /root/go/bin/gospider spider --website=iciba"
			util.ExecCommand(script)
		case 4:
			script := "docker exec -it spider /root/go/bin/gospider spider --website=googlenews"
			util.ExecCommand(script)
		default:
			break exit
		}
	}
}

// 爬取scrapGoogleNews
func scrapGoogleNews(keyword string) {

	articles := GetArticles(VenBaseUrl, keyword)
	SaveArticle(articles)

}

// 爬取每日一词
func scrapIciba() {

	url := "http://news.iciba.com/"
	selector := "body > div.screen > div.banner > div.swiper-container-place > div > div.swiper-slide.swiper-slide-0.swiper-slide-visible.swiper-slide-active > a.item.item-big > div.item-bottom"
	htmlContent, err := GetHTTPHtmlContent(url, selector, DocBodySelector)
	if err != nil {
		log.Fatal(err)
	}
	sentenceList, err := GetDataList(htmlContent, ".chinese")
	if err != nil {
		log.Fatal(err)
	}

	var todaySentence string
	sentenceList.Each(func(i int, selection *goquery.Selection) {
		todaySentence = selection.Text()
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(todaySentence)
}

// 爬取douban最佳电影
func scrapDouban() {

	movieSelector := "#content > div > div.article > ol > li"

	fmt.Println("获取页面！")
	pages := GetPages(DoubanBaseUrl)

	var movies []DoubanMovie
	for index, page := range pages {
		fmt.Printf("获取 %d 页！\n", index)
		pageContent, _ := GetHTTPHtmlContent(page.Url, movieSelector, DocBodySelector)
		pageDom, err := goquery.NewDocumentFromReader(strings.NewReader(pageContent))
		if err != nil {
			log.Fatal(err)
		}
		pageMovies := ParseMovies(pageDom)
		movies = append(movies, pageMovies...)
	}

	// save movies
	fmt.Println("保存电影记录到数据库！")
	SaveMovies(movies)
}

// GetHTTPHtmlContent 获取网站上爬取的数据
// url [string] 网址
// selector [string] 必须显示的元素
// sel [interface] 要抓取的元素
func GetHTTPHtmlContent(url string, selector string, sel interface{}) (string, error) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	//初始化参数先传一个空的数据
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	// create context
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	// 执行一个空task, 用提前创建Chrome实例
	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	//创建一个上下文超时时间为40s
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 40*time.Second)
	defer cancel()

	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector),
		chromedp.OuterHTML(sel, &htmlContent, chromedp.ByJSPath),
	)
	if err != nil {
		return "", err
	}

	return htmlContent, nil
}

// GetDataList 得到数据列表
func GetDataList(htmlContent string, selector string) (*goquery.Selection, error) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	list := dom.Find(selector)
	return list, nil
}

func init() {
	rootCmd.AddCommand(spiderCmd)
	spiderCmd.Flags().StringVarP(&webSite, "website", "t", "douban", "Choose a website")
	spiderCmd.Flags().StringVarP(&keyWord, "keyword", "k", "venezuela", "Choose a keyword to query")
}
