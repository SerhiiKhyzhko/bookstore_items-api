package config

import (
	"fmt"
	"os"

	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
)

var (
	RestyBaseUrl string
	EsHosts      string
)

func Init() {
	RestyBaseUrl = getRequiredEnv("OAUTH_API_BASE_URL")
	EsHosts = getRequiredEnv("ES_HOST_ADDRESSES")
}

func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Error(fmt.Sprintf("Critical environment variable %s is missing", key), nil)
		os.Exit(1)
	}
	return value
}
