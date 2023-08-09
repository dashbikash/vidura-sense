package robotstxtutil

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dashbikash/vidura-sense/internal/datastorage/natsio"
	"github.com/dashbikash/vidura-sense/internal/system"
	"github.com/temoto/robotstxt"
)

func getRobotsTxtCache(domainName string) string {

	return natsio.KVGet(system.Config.Data.NatsIO.KvBuckets.RobotsTxt, domainName)
}
func setRobotsTxtCache(domainName string, robotsTxt string) {
	natsio.KVPut(system.Config.Data.NatsIO.KvBuckets.RobotsTxt, domainName, robotsTxt)
}

func fetchRobotsTxtFromServer(scheme string, hostUrl string) string {

	client := &http.Client{Timeout: time.Second * 2}

	req, err := http.NewRequest("GET", scheme+"://"+hostUrl+"/robots.txt", nil)
	if err != nil {
		system.Log.Fatal(err.Error())
	}
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

func GetRobotsTxtForUrl(targetUrl string) string {
	robotsTxt := ""
	urlParsed, err := url.Parse(targetUrl)
	system.Log.Debug("Getting Robotstxt: " + urlParsed.Host)
	if err != nil {
		system.Log.Error(err.Error())
	} else {
		urlParsed.Host = strings.TrimPrefix(urlParsed.Host, "www.")
		robotsTxt = getRobotsTxtCache(urlParsed.Host)
		if len(robotsTxt) < 1 {
			robotsTxt = fetchRobotsTxtFromServer(urlParsed.Scheme, urlParsed.Host)
			setRobotsTxtCache(urlParsed.Host, robotsTxt)
		}
	}

	return robotsTxt

}

func IsAllowedUrl(targetUrl string) bool {
	robotsTxtRules := GetRobotsTxtForUrl(targetUrl)
	if robotsTxtRules == "na" {
		return true
	}
	robots, err := robotstxt.FromString(robotsTxtRules)
	if err != nil {
		system.Log.Error(err.Error())
	}

	urlParsed, err := url.Parse(targetUrl)
	system.Log.Debug("Getting Robotstxt for:" + urlParsed.Host)
	if err != nil {
		system.Log.Error(err.Error())
	}

	allow := robots.TestAgent(urlParsed.Path, system.Config.Crawler.UserAgent)

	return allow
}
