package main

import (
	"basic"
	"controller"
	"fmt"
)

type Person struct {
	name string
}

func main() {
	fmt.Println("go!")
	controller := controller.NewController()
	basic.SetConfig("mainurl", "http://www.ccse.uestc.edu.cn/")
	basic.SetConfig("req_method", "GET")
	controller.Go()

	fmt.Println("main quit!")
}
