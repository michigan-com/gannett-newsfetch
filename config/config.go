package config

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	MongoUri            string `envconfig:"mongo_uri"`
	GannettSearchApiKey string `envconfig:"gannett_search_api_key"`
	GannettAssetApiKey  string `envconfig:"gannett_asset_api_key"`
	SiteCodes           string `envconfig:"site_codes"`
	SummaryVEnv         string `envconfig:"summary_v_env"`
	GnapiDomain         string `envconfig:"gnapi_domain"`
}

type ApiConfig struct {
	GannettSearchApiKey string
	GannettAssetApiKey  string
	SiteCodes           []string
}

func GetApiConfig() (apiConfig ApiConfig, err error) {
	var env EnvConfig
	err = envconfig.Process("gannett-newsfetch.api", &env)

	apiConfig.GannettSearchApiKey = env.GannettSearchApiKey
	apiConfig.GannettAssetApiKey = env.GannettAssetApiKey
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
