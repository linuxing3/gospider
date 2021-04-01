package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/linuxing3/gospider/config"
	"github.com/linuxing3/vpsman/util"
	"github.com/spf13/cobra"
)

var (
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
		} else if webSite == "cda" {
			scrapCda(keyWord)
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
		loopMenu := []string{"iciba", "douban", "hacknews", "googlenews", "cda"}
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
		case 5:
			script := "docker exec -it spider /root/go/bin/gospider spider --website=cda"
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
	htmlContent, err := config.GetHTTPHtmlContent(url, selector, config.DocBodySelector)
	if err != nil {
		log.Fatal(err)
	}
	sentenceList, err := config.GetDataList(htmlContent, ".chinese")
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
		pageContent, _ := config.GetHTTPHtmlContent(page.Url, movieSelector, config.DocBodySelector)
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

func init() {
	rootCmd.AddCommand(spiderCmd)
	spiderCmd.Flags().StringVarP(&webSite, "website", "t", "douban", "Choose a website")
	spiderCmd.Flags().StringVarP(&keyWord, "keyword", "k", "venezuela", "Choose a keyword to query")
}
