package system

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

var config = GetConfig()

func GetLogger() *zap.Logger {

	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = config.Application.Log.Outputs
	cfg.Level = config.Application.Log.Level

	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := cfg.Build()

	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return logger
}
func GetConfig() *Config {
	configFile := "config/config.yml"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	ymlText, err := ReadTextFile(configFile)
	if err != nil {
		panic("Failed to read configuration file")
	}
	cf := &Config{}
	err = yaml.Unmarshal([]byte(ymlText), &cf)
	if err != nil {
		panic(err)
	}

	return cf
}

func ReadTextFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

type Config struct {
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
