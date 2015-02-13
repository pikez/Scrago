package analyzer

import (
	"basic"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

type GenAnalyzer interface {
	ParsePage(httpRes *http.Response) ([]string, []basic.Item)
}

type Analyzer struct {
	linklist []string
	itemlist []basic.Item
}

func NewAnalyzer() GenAnalyzer {
	return &Analyzer{
		make([]string, 0),
		make([]basic.Item, 0),
	}
}

//用于解析页面
func (self *Analyzer) ParsePage(httpRes *http.Response) ([]string, []basic.Item) {
	doc, _ := goquery.NewDocumentFromResponse(httpRes)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, exits := s.Attr("href")
		if exits {
			link = basic.CheckLink(link)
			if link != "" {
				self.linklist = append(self.linklist, link)
			}
			text := strings.TrimSpace(s.Text())
			if text != "" {
				item := make(map[string]interface{})
				item["标题"] = text
				self.itemlist = append(self.itemlist, item)
			}
		}
	})
	httpRes.Body.Close()
	return self.linklist, self.itemlist
}
