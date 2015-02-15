package middleware

import ()

type WorkPool struct {
}

func NewWorkPool() *WorkPool {
	return &WorkPool{}
}

//以num个goroutine执行函数
func (self *WorkPool) Pool(num int, work func()) {
	for w := 0; w < num; w++ {
		go work()
	}
}
