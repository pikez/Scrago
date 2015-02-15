package analyzer

import (
	"basic"
	//"fmt"
	"net/http"
)

type GenAnalyzer interface {
	Analyze(httpRes *http.Response, parser Parser) ([]string, []basic.Item)
}

type Parser func(httpRes *http.Response) ([]string, []basic.Item)

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
func (self *Analyzer) Analyze(httpRes *http.Response, parser Parser) ([]string, []basic.Item) {
	defer httpRes.Body.Close()
	if parser == nil {
		panic("xxx")
	}
	return parser(httpRes)
}
