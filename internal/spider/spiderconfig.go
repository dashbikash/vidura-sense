package spider

type SpiderConfig struct {
	IgnoreRobotTxt bool
	Proxies        []string
}

func DefaultConfig() *SpiderConfig {
	return &SpiderConfig{
		IgnoreRobotTxt: false,
		Proxies:        []string{},
	}
}
