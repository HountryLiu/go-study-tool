package go_query

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/HountryLiu/go-study-tool/model"
	"github.com/HountryLiu/go-study-tool/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"github.com/gin-gonic/gin"
	"github.com/gogs/chardet"
)

// @Tags GoQuery
// @Summary GoQuery使用
// @Description GoQuery使用
// @accept application/json
// @Produce application/json
// @Failure 200 {object} object{no=int,data=string}
// @Router /api/goquery [get]
func Index(ctx *gin.Context) {

	SpiderPages(ctx, 1, 10)
	utils.Success(ctx)
}

/*
由于 net/html 要求使用 UTF-8 编码，goquery 也是如此。我们需要保证传给 goquery 的 HTML 源字符串是 UTF-8 编码的
*/
func BodyReaderToUtf8(body io.Reader) (uft8_body io.Reader, err error) {
	content, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	charset, err := chardet.NewHtmlDetector().DetectBest(content)
	if err != nil {
		return
	}
	//body在上方已经被读过了，需要重置才能再读
	body = ioutil.NopCloser(bytes.NewBuffer(content))
	if charset.Charset != "UTF-8" {
		uft8_body, err = iconv.NewReader(body, charset.Charset, "UTF-8")
		if err != nil {
			return
		}
	} else {
		uft8_body = body
	}
	return
}

func trimStr(str string) string {
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	str = strings.Replace(str, "<br />", "", -1)
	str = strings.Replace(str, "&nbsp;", "", -1)
	return str
}

func SpiderPages(ctx *gin.Context, start, end int) {
	fmt.Printf("正在爬取 %d 到 %d 的页面\n", start, end)
	page := make(chan int)

	for i := start; i <= end; i++ {
		go func(i int) {
			err := SpiderPage(i, page)
			if err != nil {
				utils.Error(ctx, utils.InternalServerError, err)
			}
		}(i)
	}

	for i := start; i <= end; i++ {
		//channel阻塞
		fmt.Printf("第%d个页面爬取完成\n", <-page)
	}
}

func SpiderPage(i int, page chan<- int) (err error) {
	// 1. 目标网址
	/*
		http://www.xiaohua8.com/xiaohua/shuxuexiaohua_1
		http://www.xiaohua8.com/xiaohua/shuxuexiaohua_2
		http://www.xiaohua8.com/xiaohua/shuxuexiaohua_3
	*/
	baseUrl := "http://www.xiaohua8.com/xiaohua/shuxuexiaohua_"
	url := baseUrl + strconv.Itoa(i)
	fmt.Printf("正在爬取第%d个网页：%s\n", i, url)

	// 2. 将所有网站的内容抓取下来
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = fmt.Errorf("url request failed: status code(%d)", res.StatusCode)
		return
	}
	uft8_body, err := BodyReaderToUtf8(res.Body)
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(uft8_body)
	if err != nil {
		return
	}

	datas := []model.GoQueryData{}
	doc.Find(".content li").Each(func(i int, s *goquery.Selection) {
		data := model.GoQueryData{}
		//2.1. 获取标题
		title := s.Find(".titles")
		data.Title = title.Text()

		//2.2. 获取内容
		if joy_url, ok := title.Attr("href"); ok {
			content, err := SpiderOneJoy(joy_url)
			if err != nil {
				return
			}
			data.Url = joy_url
			data.Content = content
		}

		datas = append(datas, data)
	})

	// 3. 把内容存入数据库
	if err = model.DB().Create(&datas).Error; err != nil {
		return
	}

	page <- i
	return
}

func SpiderOneJoy(url string) (content string, err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = fmt.Errorf("url request failed: status code(%d)", res.StatusCode)
		return
	}
	uft8_body, err := BodyReaderToUtf8(res.Body)
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(uft8_body)
	if err != nil {
		return
	}
	content = trimStr(doc.Find("#article-content1 p").Text())
	return
}
