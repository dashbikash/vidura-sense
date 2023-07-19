package provider

import (
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func GetLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	return logger
}

func GetConfig() Config {
	ymlText, err := ReadTextFile("config.yml")
	if err != nil {
		panic(err)
	}
	config := Config{}
	err = yaml.Unmarshal([]byte(ymlText), &config)
	if err != nil {
		panic(err)
	}

	return config
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
	}
	Server struct {
		Mode string
		Port string
		Host string
	}

	Crawler struct {
		UserAgent string
		Proxies   []string
	}
}
