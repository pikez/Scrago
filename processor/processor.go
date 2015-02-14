package processor

import (
	"basic"
	"net/http"
)

type GenProcessor interface {
	DealLink(link basic.Link) (*basic.Request, bool)
	GetVurl() map[string]bool
	DealItem(item basic.Item)
}

type Processor struct {
	Vurl map[string]bool //已访问过的url字典
}

func NewProcessor() *Processor {
	return &Processor{make(map[string]bool)}
}

func (self *Processor) DealLink(link basic.Link) (*basic.Request, bool) {
	//url判重
	if _, visited := self.Vurl[link.GetLink()]; visited {
		return nil, false
	}
	self.Vurl[link.GetLink()] = true //放入字典
	httpReq, err := http.NewRequest(basic.Config.RequestMethod, link.GetLink(), nil)
	basic.Check(err)
	request := basic.NewRequest(httpReq, link.GetIndex()) //转化为构造请求
	return request, true
}

func (self *Processor) DealItem(item basic.Item) {

}

func (self *Processor) GetVurl() map[string]bool {
	return self.Vurl
}
