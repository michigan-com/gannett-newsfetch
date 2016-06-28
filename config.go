package main

import (
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MongoURI           string `envconfig:"mongo_uri"`
	GannettAPIKey      string `envconfig:"gannett_search_api_key"`
	GannettAssetAPIKey string `envconfig:"gannett_asset_api_key"`
	SiteCodes          []string
	SummaryVEnv        string `envconfig:"summary_v_env"`
	GNAPIDomain        string `envconfig:"gnapi_domain"`

	BrvtyURL     string `envconfig:"brvty_url"`
	BrvtyAPIKey  string `envconfig:"brvty_api_key"`
	BrvtyTimeout time.Duration

	SiteCodesStr   string `envconfig:"site_codes"`
	BrvtyTimeoutMs int    `envconfig:"brvty_timeout"`

	LoopInterval time.Duration
}

func ParseConfig() (Config, error) {
	var config Config
	err := envconfig.Process("newsfetch", &config)
	if err != nil {
		return config, err
	}

	config.SiteCodes = strings.Split(config.SiteCodesStr, ",")
	config.BrvtyTimeout = time.Duration(config.BrvtyTimeoutMs) * time.Millisecond

	return config, err
}
