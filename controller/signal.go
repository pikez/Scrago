package controller

import (
	"sync"
)

type GenStopSignal interface {
	sign() bool
	status() bool
	finish() bool
	ended() bool
}

func NewStopSignal() *StopSignal {
	return &StopSignal{}
}

type StopSignal struct {
	rwmutex sync.RWMutex //读写锁
	flag    bool         //停止信号是否发出的标志
	stop    bool         //模块是否停止完全的标志
}

func (self *StopSignal) sign() bool {
	self.rwmutex.Lock()
	defer self.rwmutex.Unlock()
	if self.flag {
		return false
	}
	self.flag = true
	return true
}

func (self *StopSignal) status() bool {
	return self.flag
}

func (self *StopSignal) finish() bool {
	self.rwmutex.Lock()
	defer self.rwmutex.Unlock()
	if self.stop {
		return false
	}
	self.stop = true
	return true
}

func (self *StopSignal) ended() bool {
	return self.stop
}
