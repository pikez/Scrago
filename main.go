package main

import (
	"basic"
	"controller"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

func main() {
	//创建一个控制器，这里有4个必须给与的参数：
	//爬取的初始url，爬取深度，解析函数，储存函数
	//basic.InitConfig() //配置Config需要先初始化
	//basic.Config.HttpHeader["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1; WOW64) " +
	//	"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.93 Safari/537.36"
	controller := controller.NewController("http://z3lion.com/", 1, Parser, Store)
	//启动爬虫
	controller.Go()
}

//解析函数定义，声明必须相同，下同
func Parser(httpRes *http.Response) ([]string, []basic.Item) {
	//两个需要返回的列表
	linklist := make([]string, 0)     //下一步需要请求的链接
	itemlist := make([]basic.Item, 0) //保存的数据类型为 map[string]interface{}

	//自定义部分
	//抓取所有链接
	doc, _ := goquery.NewDocumentFromResponse(httpRes)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, exits := s.Attr("href")
		if exits {
			link = basic.CheckLink(link)
			if link != "" {
				linklist = append(linklist, link)
			}
		}
	})
	//保存每个页面的标题
	title := strings.TrimSpace(doc.Find("head title").Text())
	if title != "" {
		item := make(map[string]interface{})
		item["标题"] = title
		itemlist = append(itemlist, item)
	}

	return linklist, itemlist
}

//储存函数定义
func Store(item basic.Item) {
	//这里只打印抓取的数据
	fmt.Println(item)
}
