package basic

import ()

type config struct {
	Name             string
	StartUrl         string
	RequestMethod    string
	DownloaderNumber int
	AnalyzerNumber   int
	ProcessorNumber  int
	ReqChanLength    int
	ResChanLength    int
	LinkChanLength   int
	ItemChanLength   int
}

var Config config

func InitConfig() {
	if Config.StartUrl == "" {
		panic("StartUrl Can not be empty!  ")
	}
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

}
