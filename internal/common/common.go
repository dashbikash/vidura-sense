package common

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var config = GetConfig()

func GetLogger() *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	if config.Application.Log.Output == "file" {
		logFile := config.Application.Log.Dir + "/" + "app.log"
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.Error("Failed to log to file, using default stdout")
		} else {
			log.Out = file
		}
	}
	log.Level = logrus.DebugLevel
	return log
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
			Level  string
			Output string
			Dir    string
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
