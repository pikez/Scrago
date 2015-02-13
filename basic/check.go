package basic

import (
	"net/url"
	"strings"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func CheckBaseurl(Url string) string {
	u, _ := url.Parse(Url)
	if u.Scheme == "" {
		Url = "http://" + Url
	}
	if flag := strings.HasSuffix(Url, "/"); flag != true {
		Url = Url + "/"
	}
	return Url
}

func CheckLink(link string) string {
	u, _ := url.Parse(link)
	if u.Scheme != "" {
		return ""
	}
	if u.Scheme == "http" || u.Scheme == "https" {
		return link
	}
	if flag := strings.HasPrefix(link, Config["mainurl"]); flag != true {
		link = strings.Join([]string{Config["mainurl"], link}, "")
		return link
	}
	return ""
}
