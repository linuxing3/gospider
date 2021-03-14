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

// spiderCmd represents the spider command
var spiderCmd = &cobra.Command{
	Use:   "spider",
	Short: "A brief description of your command",
	Long: ` About usage of using spider. For example: 
Trojan is a CLI Command for Go that empowers proxy.
to quickly create a web tunnel.`,
	Run: func(cmd *cobra.Command, args []string) {
		spiderMenu()
	},
}

// TrojanMenu 控制TrojanMenu
func spiderMenu() {
exit:
	for {
		fmt.Println()
		fmt.Print(util.Cyan("请选择"))
		fmt.Println()
		loopMenu := []string{"iciba", "douban", "hacknews", "other"}
		choice := util.LoopInput("回车退出", loopMenu, false)
		switch choice {
		case 1:
			scrapIciba()
		case 2:
			scrapDouban()
		case 3:
			scrapIciba()
		case 4:
			scrapIciba()
		default:
			break exit
		}
	}
}

func scrapIciba() {

	url := "http://news.iciba.com/"
	selector := "body > div.screen > div.banner > div.swiper-container-place > div > div.swiper-slide.swiper-slide-0.swiper-slide-visible.swiper-slide-active > a.item.item-big > div.item-bottom"
	sel := "document.querySelector('body')"
	htmlContent, err := GetHTTPHtmlContent(url, selector, sel)
	if err != nil {
		log.Fatal(err)
	}
	todaySentence, err := GetSpecialData(htmlContent, ".chinese")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(todaySentence)
}

func scrapDouban() {
	var movies []DoubanMovie

	pages := GetPages(DoubanBaseUrl)
	for _, page := range pages {
		doc, err := goquery.NewDocument(strings.Join([]string{DoubanBaseUrl, page.Url}, ""))
		if err != nil {
			log.Println(err)
		}

		movies = append(movies, ParseMovies(doc)...)
	}
    log.Println(movies)
}


// GetHTTPHtmlContent 获取网站上爬取的数据
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

// GetSpecialData 得到具体的数据
func GetSpecialData(htmlContent string, selector string) (string, error) {

    list, err := GetDataList(htmlContent, selector) 
	if err != nil {
		return "", err
	}
    
    var str string
    list.Each(func(i int, selection *goquery.Selection) {
		str = selection.Text()
	})
	return str, nil
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// spiderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// spiderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
