package system

import (
	"os"

	"github.com/dashbikash/vidura-sense/internal/helper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var Config *SystemConfig

type SystemConfig struct {
	Application struct {
		ID      string
		Version string
		Log     struct {
			Level   zap.AtomicLevel
			Outputs []string
		}
	}
	Server struct {
		Mode string
		Port string
		Host string
	}

	Crawler struct {
		UserAgent string `yaml:"user-agent"`
		Proxies   []string
	}
	Data struct {
		Mongo struct {
			MongoUrl    string `yaml:"mongo-url"`
			Database    string
			Collections struct {
				Htmlpages string `yaml:"htmlpages"`
				Feeditems string `yaml:"feeditems"`
			}
		}
		Redis struct {
			RedisUrl string `yaml:"redis-url"`
			Branches struct {
				RobotsTxt struct {
					Name string
					Ttl  int
				} `yaml:"robots-txt"`
			}
		}
		NatsIO struct {
			NatsUrl   string `yaml:"nats-url"`
			KvBuckets struct {
				RobotsTxt string `yaml:"robots-txt"`
			} `yaml:"kv-buckets"`
		}
	}
}

func setConfig() {

	configFile := "config/config.yml"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	ymlText, err := helper.ReadTextFile(configFile)
	if err != nil {
		panic("Failed to read configuration file")
	}
	Config = &SystemConfig{}
	err = yaml.Unmarshal([]byte(ymlText), &Config)
	if err != nil {
		panic(err)
	}
}
