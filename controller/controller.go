package controller

import (
	"analyzer"
	"basic"
	"downloader"
	"fmt"
	"middleware"
	"net/http"
	"processor"
	"sync"
	"time"
)

var Maxindex uint32 = 1
var wg sync.WaitGroup //全局wait锁

type Controller struct {
	Downloader downloader.GenDownloader //下载器
	Analyzer   analyzer.GenAnalyzer     //分析器
	Processor  processor.GenProcessor   //处理器
	Channel    *middleware.Channel      //管道
	WorkPool   *middleware.WorkPool     //工作池
	StopSignal *StopSignal              //停止信号
}

func NewController() *Controller {
	return &Controller{}
}

func (ctrl *Controller) Go() {
	ctrl.Downloader = downloader.NewDownloader() //初始化各组件，下同
	ctrl.Analyzer = analyzer.NewAnalyzer()
	ctrl.Processor = processor.NewProcessor()
	ctrl.Channel = middleware.NewChannel()
	ctrl.WorkPool = middleware.NewWorkPool()
	ctrl.StopSignal = NewStopSignal()

	//准备第一次请求
	primaryreq, err := http.NewRequest(basic.Config["method"], basic.Config["mainurl"], nil)
	basic.Check(err)
	basereq := basic.NewRequest(primaryreq, 0)
	ctrl.Channel.ReqChan() <- *basereq

	//利用goroutine使三个组件同时工作，启动监视器监视
	wg.Add(4)
	go ctrl.DownloaderManager()
	go ctrl.AnalyzerManager()
	go ctrl.ProcessorManager()
	go ctrl.Monitors()

	wg.Wait()
	fmt.Println(len(ctrl.Processor.GetVurl()))
}

func (ctrl *Controller) DownloaderManager() {
	defer wg.Done()
	dwg := new(sync.WaitGroup)
	dwg.Add(5)
	//工作池机制
	ctrl.WorkPool.Pool(5, func() {
		for req := range ctrl.Channel.ReqChan() {
			res := ctrl.Downloader.Download(&req) //res为构造请求类型
			if res != nil {
				ctrl.Channel.ResChan() <- *res //放入响应通道
			}
		}
		dwg.Done()
	})
	dwg.Wait()
	//请求通道关闭后关闭响应通道
	close(ctrl.Channel.ResChan())
	fmt.Println("download quit!")
}

func (ctrl *Controller) AnalyzerManager() {
	defer wg.Done()
	awg := new(sync.WaitGroup)
	awg.Add(2)
	ctrl.WorkPool.Pool(2, func() {
		for res := range ctrl.Channel.ResChan() {
			Links, Items := ctrl.Analyzer.ParsePage(res.GetRes()) //解析函数解析html页面
			//将item放入通道传至持久储存函数
			for _, item := range Items {
				ctrl.Channel.ItemChan() <- item
			}
			//如果停止信号发出，不再向链接通道传输数据
			if ctrl.StopSignal.status() {
				continue
			}
			//发送至链接通道后续处理
			for _, link := range Links {
				ctrl.Channel.LinkChan() <- basic.NewLinks(link, res.GetIndex()+1)
			}
		}
		awg.Done()
	})
	awg.Wait()
	//同上，发送完毕后关闭
	close(ctrl.Channel.LinkChan())
	close(ctrl.Channel.ItemChan())
	//运行流程的最后一步，所以发送已结束信号
	ctrl.StopSignal.finish()
	fmt.Println("Analyzer quit!")
}

func (ctrl *Controller) ProcessorManager() {
	defer wg.Done()
	pwg := new(sync.WaitGroup)
	pwg.Add(4)
	ctrl.WorkPool.Pool(2, func() {
		for link := range ctrl.Channel.LinkChan() {
			//flag为判重标志
			req, flag := ctrl.Processor.DealLink(link)
			if !flag {
				continue
			}
			//判断深度（索引）
			if req.GetIndex() <= Maxindex {
				ctrl.Channel.ReqChan() <- *req
			} else {
				//避免重复关闭通道
				if ctrl.StopSignal.status() {
					continue
				}
				close(ctrl.Channel.ReqChan())
				ctrl.StopSignal.sign()
			}
		}
		pwg.Done()
	})
	ctrl.WorkPool.Pool(2, func() {
		for item := range ctrl.Channel.ItemChan() {
			ctrl.Processor.DealItem(item)
		}
		pwg.Done()
	})
	pwg.Wait()
	fmt.Println("Processor quit!")
}

func (ctrl *Controller) Monitors() {
	defer wg.Done()
	begin := time.Now()
	for {
		fmt.Println(time.Now().Sub(begin))
		fmt.Println("reqchan length", len(ctrl.Channel.ReqChan()))
		fmt.Println("reschan length :", len(ctrl.Channel.ResChan()))
		fmt.Println("linkchan length", len(ctrl.Channel.LinkChan()))
		fmt.Println("itemchan length", len(ctrl.Channel.ItemChan()))
		time.Sleep(time.Second)
		if ctrl.StopSignal.ended() {
		}
	}

}
