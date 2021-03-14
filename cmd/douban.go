package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/linuxing3/gospider/prisma/db"
)

type DoubanMovie struct {
	Title    string
	Subtitle string
	Other    string
	Desc     string
	Year     string
	Area     string
	Tag      string
	Star     string
	Comment  string
	Quote    string
}

type Page struct {
	Page int
	Url  string
}

var (
	DoubanBaseUrl = "https://movie.douban.com/top250"
	DoubanTopPageSelector = "#content > div > div.article > div.paginator > a"
)

// GetPages 获取分页
func GetPages(url string) (pages []Page) {
	htmlContent, err := GetHTTPHtmlContent(DoubanBaseUrl, DoubanTopPageSelector, DocBodySelector)
	if err != nil {
		log.Fatal(err)
	}

	pageList, err := GetDataList(htmlContent, DoubanTopPageSelector)
	if err != nil {
		log.Fatal("No list")
	}
	pageList.Each(func(i int, selection *goquery.Selection) {
		title, _ := strconv.Atoi(selection.Text())
		url, _ := selection.Attr("href")
		pageUrl := strings.Join([]string{DoubanBaseUrl, url}, "")
		pages = append(pages, Page{
			Page: title,
			Url:  pageUrl,
		})
	})
	return pages
}

// ParseMovies 在每一个页面上分析电影数据
func ParseMovies(doc *goquery.Document) (movies []DoubanMovie) {
	movieSelector := "#content > div > div.article > ol > li"
  doc.Find(movieSelector).Each(func(i int, s *goquery.Selection) {

		fmt.Printf("获取第 %d 个电影\n", i)

		title := s.Find(".hd > a > span").Eq(0).Text()

		subtitle := s.Find(".hd > a > span").Eq(1).Text()
		subtitle = strings.TrimLeft(subtitle, "  / ")

		other := s.Find(".hd > a > span").Eq(2).Text()
		other = strings.TrimLeft(other, "  / ")

		desc := strings.TrimSpace(s.Find(".bd > p").Eq(0).Text())
		DescInfo := strings.Split(desc, "\n")
		desc = DescInfo[0]

		movieDesc := strings.Split(DescInfo[1], "/")
		year := strings.TrimSpace(movieDesc[0])
		area := strings.TrimSpace(movieDesc[1])
		tag := strings.TrimSpace(movieDesc[2])

		star := s.Find(".bd > .star > .rating_num").Text()

		comment := strings.TrimSpace(s.Find(".bd > .star > span").Eq(3).Text())
		compile := regexp.MustCompile("[0-9]")
		comment = strings.Join(compile.FindAllString(comment, -1), "")

		quote := s.Find(".quote > .inq").Text()

		movie := DoubanMovie{
			Title:    title,
			Subtitle: subtitle,
			Other:    other,
			Desc:     desc,
			Year:     year,
			Area:     area,
			Tag:      tag,
			Star:     star,
			Comment:  comment,
			Quote:    quote,
		}

		// log.Printf("i: %d, movie: %v", i, movie)

		movies = append(movies, movie)
	})

	return movies
}

// SaveMovies 保存电影记录到数据库
func SaveMovies(movies []DoubanMovie)  {
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
	for _, movie := range movies {
		_, err := client.Movies.CreateOne(
			db.Movies.Title.Set(movie.Title),
			db.Movies.Subtitle.Set(movie.Subtitle),
			db.Movies.Desc.Set(movie.Desc),
			db.Movies.Area.Set(movie.Area),
			db.Movies.Year.Set(movie.Year),
			db.Movies.Tag.Set(movie.Tag),
			db.Movies.Star.Set(movie.Star),
		).Exec(ctx)
    if err != nil {
			fmt.Println(err)
    }
	}
	
}

// ExampleScrape 测试抓取网页
func ExampleScrape() {
  // Request the HTML page.
  res, err := http.Get("http://metalsucks.net")
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()
  if res.StatusCode != 200 {
    log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  }

  // Load the HTML document
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }

	fmt.Println("Example scrapy")
  // Find the review items
  doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
    // For each item found, get the band and title
    band := s.Find("a").Text()
    title := s.Find("i").Text()
    fmt.Printf("Review %d: %s - %s\n", i, band, title)
  })
}