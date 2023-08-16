package robotstxtutil

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/dashbikash/vidura-sense/internal/datastorage/natsio"
	"github.com/dashbikash/vidura-sense/internal/system"
	"github.com/temoto/robotstxt"
)

type RobotsTxt struct {
	proxy     *url.URL
	value     string
	targetUrl string
}

func (rt *RobotsTxt) getRobotsTxtCache(domainName string) string {

	return natsio.KVGet(system.Config.Data.NatsIO.KvBuckets.RobotsTxt, domainName)
}
func (rt *RobotsTxt) setRobotsTxtCache(domainName string, robotsTxt string) {
	natsio.KVPut(system.Config.Data.NatsIO.KvBuckets.RobotsTxt, domainName, robotsTxt)
}

func (rt *RobotsTxt) fetchRobotsTxtFromServer(scheme string, hostUrl string) string {

	client := &http.Client{Timeout: time.Second * 5}

	req, err := http.NewRequest("GET", scheme+"://"+hostUrl+"/robots.txt", nil)
	if err != nil {
		system.Log.Fatal(err.Error())
	}
	client.Transport = &http.Transport{Proxy: http.ProxyURL(rt.proxy)}
	req.Header.Set("User-Agent", system.Config.Crawler.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		system.Log.Error(err.Error())
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		system.Log.Error(err.Error())
		return ""
	}
	defer func() {
		resp.Body.Close()
		req = nil
		client = nil
		err = nil
	}()
	if resp.StatusCode == http.StatusOK {
		robotsVal := string(body)
		return robotsVal
	} else if resp.StatusCode == http.StatusNotFound {
		return "na"
	}
	return ""

}

func RobotsTxtFor(targetUrl string, proxy *url.URL) *RobotsTxt {
	rt := &RobotsTxt{targetUrl: targetUrl, proxy: proxy}
	urlParsed, err := url.Parse(targetUrl)
	system.Log.Debug("Getting robots.txt: for" + urlParsed.Host)
	if err != nil {
		system.Log.Error(err.Error())
	} else {
		rt.value = rt.getRobotsTxtCache(urlParsed.Host)
		if len(rt.value) < 1 {
			rt.value = rt.fetchRobotsTxtFromServer(urlParsed.Scheme, urlParsed.Host)
			rt.setRobotsTxtCache(urlParsed.Host, rt.value)
		}
	}

	return rt

}

func (rt *RobotsTxt) IsAllowed() bool {

	if rt.value == "na" {
		return true
	}
	robots, err := robotstxt.FromString(rt.value)
	if err != nil {
		system.Log.Error(err.Error())
	}

	urlParsed, err := url.Parse(rt.targetUrl)
	if err != nil {
		system.Log.Error(err.Error())
	}

	allow := robots.TestAgent(urlParsed.Path, system.Config.Crawler.UserAgent)

	return allow
}
