package basic

import ()

type Link struct {
	link  string
	index uint32
}

func NewLinks(link string, index uint32) Link {
	return Link{link, index}
}

func (self *Link) GetLink() string {
	return self.link
}

func (self *Link) GetIndex() uint32 {
	return self.index
}
