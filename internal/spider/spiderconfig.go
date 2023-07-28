package spider

import "context"

type SpiderConfig struct {
	IgnoreRobotTxt bool
	Proxies        []string
	Ctx            context.Context
}

func DefaultConfig() *SpiderConfig {
	return &SpiderConfig{
		IgnoreRobotTxt: false,
		Proxies:        []string{},
		Ctx:            context.Background(),
	}
}
