package system

import (
	"os"

	"github.com/dashbikash/vidura-sense/internal/helper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

var Config = getConfig()
var Logger = getLogger()

func getLogger() *zap.Logger {

	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = Config.Application.Log.Outputs
	cfg.Level = Config.Application.Log.Level

	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := cfg.Build()

	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return logger
}
func getConfig() *SystemConfig {
	configFile := "config/config.yml"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	ymlText, err := helper.ReadTextFile(configFile)
	if err != nil {
		panic("Failed to read configuration file")
	}
	cf := &SystemConfig{}
	err = yaml.Unmarshal([]byte(ymlText), &cf)
	if err != nil {
		panic(err)
	}

	return cf
}

type SystemConfig struct {
	Application struct {
		Name    string
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
	}
}
