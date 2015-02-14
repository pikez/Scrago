package middleware

import (
	"basic"
)

type Channel struct {
	reqchan  chan basic.Request
	reschan  chan basic.Response
	linkchan chan basic.Link
	itemchan chan basic.Item
}

func NewChannel() *Channel {
	return &Channel{
		make(chan basic.Request, basic.Config.ReqChanLength),
		make(chan basic.Response, basic.Config.ResChanLength),
		make(chan basic.Link, basic.Config.LinkChanLength),
		make(chan basic.Item, basic.Config.ItemChanLength),
	}
}

func (self *Channel) ReqChan() chan basic.Request {
	return self.reqchan
}

func (self *Channel) ResChan() chan basic.Response {
	return self.reschan
}

func (self *Channel) LinkChan() chan basic.Link {
	return self.linkchan
}

func (self *Channel) ItemChan() chan basic.Item {
	return self.itemchan
}
