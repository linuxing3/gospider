package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/linuxing3/gospider/config"
	"github.com/linuxing3/gospider/prisma/db"
)

type GoogleNewsArticle struct {
	Title    string
	Url      string
}

var (
	VenBaseUrl = "https://news.google.com/search?q=venezuela&hl=es-419&gl=VE&ceid=VE%3Aes-419"
	VenAritcleSelector = "article > h3 > a"
)

// GetPages 获取分页
func GetArticles(url string, keyword string) (articles []GoogleNewsArticle) {

	VenBaseUrl = strings.Replace(VenBaseUrl, "venezuela", keyword, 1)
	htmlContent, err := config.GetHTTPHtmlContent(VenBaseUrl, VenAritcleSelector, config.DocBodySelector)
	if err != nil {
		log.Fatal(err)
	}

	articleList, err := config.GetDataList(htmlContent, VenAritcleSelector)
	if err != nil {
		log.Fatal("No list")
	}
	articleList.Each(func(i int, selection *goquery.Selection) {
		url, _ := selection.Attr("href")
		title := selection.Text()
		pageUrl := strings.Join([]string{"https://news.google.com/", url}, "")
		articles = append(articles, GoogleNewsArticle{
			Title: title,
			Url:  pageUrl,
		})
		fmt.Println(title)
	})
	return articles
}

func SaveArticle(articles []GoogleNewsArticle)  {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()
	for _, article := range articles {
		_, err := client.Articles.CreateOne(
			db.Articles.Title.Set(article.Title),
			db.Articles.URL.Set(article.Url),
		).Exec(ctx)
    if err != nil {
			fmt.Println(err)
    }
	}
	
}

func GetArticleDetail(url string) {
	htmlContent, err := config.GetHTTPHtmlContent(url, "", config.DocBodySelector)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(htmlContent)
}