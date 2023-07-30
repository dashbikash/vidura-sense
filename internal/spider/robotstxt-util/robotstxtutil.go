package robotstxtutil

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	redisstore "github.com/dashbikash/vidura-sense/internal/datastorage/redis-store"
	"github.com/dashbikash/vidura-sense/internal/system"
	"github.com/temoto/robotstxt"
)

func getRobotsTxtCache(domainName string) string {

	return redisstore.DefaultClient().GetString(system.Config.Data.Redis.Branches.RobotsTxt.Name+":"+domainName, "")
}
func setRobotsTxtCache(domainName string, robotsTxt string) bool {
	return redisstore.DefaultClient().SetString(system.Config.Data.Redis.Branches.RobotsTxt.Name+":"+domainName, robotsTxt, time.Hour*time.Duration(system.Config.Data.Redis.Branches.RobotsTxt.Ttl))
}

func fetchRobotsTxtFromServer(hostUrl string) string {

	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://"+hostUrl+"/robots.txt", nil)
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
	if resp.StatusCode == 200 {
		robotsVal := string(body)

		return robotsVal
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
			robotsTxt = fetchRobotsTxtFromServer(urlParsed.Host)
			setRobotsTxtCache(urlParsed.Host, robotsTxt)
		}
	}

	return robotsTxt

}

func IsAllowedUrl(targetUrl string) bool {
	robotsTxtRules := GetRobotsTxtForUrl(targetUrl)
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
