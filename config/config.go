package config

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	MongoUri      string `envconfig:"mongo_uri"`
	GannettApiKey string `envconfig:"gannett_api_key"`
	SiteCodes     string `envconfig:"site_codes"`
	SummaryVEnv   string `envconfig:"summary_v_env"`
	GnapiDomain   string `envconfig:"gnapi_domain"`
	BrvtyURL      string `envconfig:"brvty_url"`
	BrvtyAPIKey   string `envconfig:"brvty_api_key"`
}

type ApiConfig struct {
	GannettApiKey string
	SiteCodes     []string
}

func GetApiConfig() (apiConfig ApiConfig, err error) {
	var env EnvConfig
	err = envconfig.Process("gannett-newsfetch.api", &env)

	apiConfig.GannettApiKey = env.GannettApiKey
	apiConfig.SiteCodes = strings.Split(env.SiteCodes, ",")

	return apiConfig, err
}

/*
	Get the current configuration from environment variables
*/
func GetEnv() (env EnvConfig, err error) {
	err = envconfig.Process("gannett-newsfetch.global", &env)
	return env, err
}
