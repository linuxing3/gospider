package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

type Course struct {
	Title    string
	Url      string
	Progress string
	Cid      string
	v        []string
}

var (
	CdaBaseURL            = "https://e-cda.cn/"
	CdaLoginFormLoginBtn  = "form > input.login_btn "
	CdaCourseRow          = "#module_0  .hoz_course_row" // rows
	CdaCourseSelector     = "h2.hoz_course_name  a"      // course name
	CdaCourseProgressBar  = "span.h_pro_percent"         // course progress %
	CdaChooseCourseBtn    = ".rt.btn_group > a"          // course video link
	CdaChooseVideoConfirm = ".continue > .user_choise"   // confirm playvideo
)

// LoginCda 登录CDA网址
func LoginCda(host string) chromedp.Tasks {

	return chromedp.Tasks{
		chromedp.Navigate(host),
		chromedp.WaitVisible(`.banner_btn`, chromedp.NodeVisible),
		chromedp.Click(".banner_btn", chromedp.NodeVisible),
		// 1. 登录
		chromedp.WaitVisible(`#username`, chromedp.ByID),
		chromedp.WaitVisible(`#pwd`, chromedp.ByID),
		chromedp.SetValue(`#username`, "xing_wenju@mfa.gov.cn", chromedp.ByID),
		chromedp.SetValue(`#pwd`, "Abcd1234", chromedp.ByID),
		chromedp.Click("input.login_btn", chromedp.NodeVisible),
	}
}

// ShowCdaCoursesList  抓取课程列表
func ShowCdaCoursesList() chromedp.Tasks {

	return chromedp.Tasks{
		// 2. 打开课程列表
		chromedp.WaitVisible(CdaCourseRow, chromedp.NodeVisible),
		chromedp.Click(CdaCourseSelector, chromedp.NodeVisible),
	}
}

// PlayCdaCourseVideo 打开课程视频网页并播放
func PlayCdaCourseVideo() chromedp.Tasks {

	return chromedp.Tasks{

		// 3. 打开视频网页
		chromedp.Click(CdaChooseCourseBtn, chromedp.NodeVisible),

		// 4. 视频播放
		chromedp.WaitVisible("video", chromedp.NodeVisible),
		chromedp.Click(CdaChooseVideoConfirm, chromedp.NodeVisible),

		// 模拟超时
		chromedp.WaitVisible("test", chromedp.NodeVisible),
	}
}

// GetCdaCoursesWithDetails 不打开视频网页,而是抓取当前页面的课程列表
func GetCdaCoursesWithDetails(cid string) (courses []Course) {
	classid := fmt.Sprintf("cid=%s", cid)
	myclass := strings.Join([]string{CdaBaseURL, "/student/class_detail_course.do?", classid, "&elective_type=1&menu=myclass&tab_index=0"}, "")

	htmlContent, err := GetHTTPHtmlContent(myclass, CdaCourseRow, DocBodySelector)
	if err != nil {
		log.Fatal(err)
	}

	courseList, err := GetDataList(htmlContent, CdaCourseRow)
	if err != nil {
		log.Fatal(err)
	}

	courseList.Each(func(i int, s *goquery.Selection) {
		// 查找课程具体信息
		item := s.Find("ul h2.hoz_course_name a")
		url, _ := item.Attr("href")
		title := item.Text()
		progress := s.Find("h_pro_percent").Text()

		var cid string

		cid, _ = s.Find(".hover_btn").Attr("onclick")
		cid = strings.ReplaceAll(cid, "addUrl(", "")
		cid = strings.ReplaceAll(cid, ")", "")

		pageURL := strings.Join([]string{CdaBaseURL, url}, "")
		courses = append(courses, Course{
			Title:    title,
			Url:      pageURL,
			Progress: progress,
			Cid:      cid,
		})
		fmt.Println(title)
	})
	return courses
}

// GetCdaCourseVideo 根据课程列表, 抓取视频下载地址
func GetCdaCourseVideo(courses []Course) []Course{
	for _, course := range courses {
		videoPageURL := strings.Join([]string{"https://e-cda.cn/portal/play.do?menu=course&id=", course.Cid}, "")
		htmlContent, _ := GetHTTPHtmlContent(videoPageURL, "video", DocBodySelector)
		mediaList, _ := GetDataList(htmlContent, "video")
		mediaList.Each(func(i int, s *goquery.Selection) {
			url, _ := s.Attr("src")
			videoDownloadLink := strings.Join([]string{"https://cdn.gwypx.com.cn/course/n", course.Cid, "/", url}, "")
			fmt.Println(videoDownloadLink)
			course.v = append(course.v, videoDownloadLink)
		})
	}
	return courses
}
