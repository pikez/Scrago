package basic

import ()

type config struct {
	flag             bool              //配置是否初始化过的标志
	Name             string            //爬虫名:)
	StartUrl         string            //初始Url，会从创建控制器是给与的参数添加
	RequestMethod    string            //http请求的方法
	HttpHeader       map[string]string //http请求的header
	DownloaderNumber int               //下载器数目
	AnalyzerNumber   int               //分析器数目
	ProcessorNumber  int               //处理器数目
	ReqChanLength    int               //请求通道长度
	ResChanLength    int               //响应通道长度
	LinkChanLength   int               //链接通道长度
	ItemChanLength   int               //数据通道长度
}

var Config *config = new(config)

func InitConfig() {
	if flag {
		return
	}
	Config.HttpHeader = make(map[string]string)
	if Config.Name == "" {
		Config.Name = "scrago"
	}
	if Config.RequestMethod == "" {
		Config.RequestMethod = "GET"
	}
	if Config.DownloaderNumber == 0 {
		Config.DownloaderNumber = 5
	}
	if Config.AnalyzerNumber == 0 {
		Config.AnalyzerNumber = 2
	}
	if Config.ProcessorNumber == 0 {
		Config.ProcessorNumber = 3
	}
	if Config.ReqChanLength == 0 {
		Config.ReqChanLength = 500
	}
	if Config.ResChanLength == 0 {
		Config.ResChanLength = 200
	}
	if Config.LinkChanLength == 0 {
		Config.LinkChanLength = 1000
	}
	if Config.ItemChanLength == 0 {
		Config.ItemChanLength = 200
	}
	flag = true
}
