package main

import (
	"analyzer"
	"basic"
	"controller"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

func main() {
	analyzer.AddParse(ParsePage)

	controller := controller.NewController("http://www.ccse.uestc.edu.cn/", 1)
	controller.Go()
}

func ParsePage(httpRes *http.Response) ([]string, []basic.Item) {
	defer httpRes.Body.Close()
	linklist := make([]string, 0)
	itemlist := make([]basic.Item, 0)

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
	title := strings.TrimSpace(doc.Find("head title").Text())
	if title != "" {
		item := make(map[string]interface{})
		item["标题"] = title
		itemlist = append(itemlist, item)
	}

	return linklist, itemlist
}
