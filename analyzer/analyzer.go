package analyzer

import (
	"basic"
	"fmt"
	"net/http"
)

type GenAnalyzer interface {
	Analyze(httpRes *http.Response) ([]string, []basic.Item)
}

type Parser func(httpRes *http.Response) ([]string, []basic.Item)

var parser Parser

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

func AddParse(parser Parser) {
	fmt.Println(parser)
	parser = parser
}

//用于解析页面

func (self *Analyzer) Analyze(httpRes *http.Response) ([]string, []basic.Item) {
	fmt.Println(parser)
	if parser == nil {
		panic("xxx")
	}
	self.linklist, self.itemlist = parser(httpRes)
	return self.linklist, self.itemlist
}
