package controller

import (
	"analyzer"
	"basic"
	"downloader"
	//"fmt"
	"middleware"
	"net/http"
	"processor"
	"sync"
	"time"
)

var wg sync.WaitGroup //全局wait锁

var logger basic.Logger = basic.NewSimpleLogger() // 日志记录器

type Controller struct {
	Downloader downloader.GenDownloader //下载器
	Analyzer   analyzer.GenAnalyzer     //分析器
	Processor  processor.GenProcessor   //处理器
	Channel    *middleware.Channel      //管道
	WorkPool   *middleware.WorkPool     //工作池
	StopSignal *StopSignal              //停止信号
	StartUrl   string                   //初始爬行Url
	Depth      uint32                   //爬行深度
	Parser     analyzer.Parser          //解析页面函数
	Store      processor.Store
}

func NewController(StartUrl string, Depth uint32, Parser analyzer.Parser, Store processor.Store) *Controller {
	return &Controller{StartUrl: StartUrl, Depth: Depth, Parser: Parser, Store: Store}
}

func (ctrl *Controller) Go() {
	basic.Config.StartUrl = ctrl.StartUrl
	basic.InitConfig()
	ctrl.Downloader = downloader.NewDownloader() //初始化各组件，下同
	ctrl.Analyzer = analyzer.NewAnalyzer()
	ctrl.Processor = processor.NewProcessor()
	ctrl.Channel = middleware.NewChannel()
	ctrl.WorkPool = middleware.NewWorkPool()
	ctrl.StopSignal = NewStopSignal()

	//准备第一次请求
	primaryreq, err := http.NewRequest(basic.Config.RequestMethod, basic.Config.StartUrl, nil)
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
	//fmt.Println(len(ctrl.Processor.GetVurl()))
	//fmt.Println(ctrl.Processor.GetVurl())
}

func (ctrl *Controller) DownloaderManager() {
	defer wg.Done()
	dwg := new(sync.WaitGroup)
	dwg.Add(basic.Config.DownloaderNumber)
	//工作池机制
	ctrl.WorkPool.Pool(basic.Config.DownloaderNumber, func() {
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
	logger.Infoln("download quit!")
}

func (ctrl *Controller) AnalyzerManager() {
	defer wg.Done()
	awg := new(sync.WaitGroup)
	awg.Add(basic.Config.AnalyzerNumber)
	ctrl.WorkPool.Pool(basic.Config.AnalyzerNumber, func() {
		for res := range ctrl.Channel.ResChan() {
			Links, Items := ctrl.Analyzer.Analyze(res.GetRes(), ctrl.Parser) //解析函数解析html页面
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
	//运行流程的最后一步，所以发送已结束信号,为监视器提供信号
	ctrl.StopSignal.finish()
	logger.Infoln("Analyzer quit!")
}

func (ctrl *Controller) ProcessorManager() {
	defer wg.Done()
	pwg := new(sync.WaitGroup)
	pwg.Add(basic.Config.ProcessorNumber)
	ctrl.WorkPool.Pool(1, func() { //分析link速度很快，基本上不会阻塞，只需分配一个goroutine
		for link := range ctrl.Channel.LinkChan() {
			//flag为判重标志
			req, flag := ctrl.Processor.DealLink(link)
			if !flag {
				continue
			}
			//判断深度（索引）
			if req.GetIndex() <= ctrl.Depth {
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
	ctrl.WorkPool.Pool(basic.Config.ProcessorNumber-1, func() {
		for item := range ctrl.Channel.ItemChan() {
			ctrl.Processor.DealItem(item, ctrl.Store)
		}
		pwg.Done()
	})
	pwg.Wait()
	logger.Infoln("Processor quit!")
}

func (ctrl *Controller) Monitors() {
	defer wg.Done()
	logger.Infoln("Spider start! ")
	for {
		logger.Infof("Spider Status: \n"+
			"reqchan length: %d\n"+
			"reschan length: %d\n"+
			"linkchan length: %d\n"+
			"itemchan length: %d\n",
			len(ctrl.Channel.ReqChan()), len(ctrl.Channel.ResChan()), len(ctrl.Channel.LinkChan()), len(ctrl.Channel.ItemChan()))
		time.Sleep(time.Second * 3)
		if ctrl.StopSignal.ended() {
			logger.Infof("Spider is Stoped!")
			break
		}
	}
}
