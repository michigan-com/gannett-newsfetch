package main

import (
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MongoUri      string `envconfig:"mongo_uri"`
	GannettApiKey string `envconfig:"gannett_api_key"`
	SiteCodes     []string
	SummaryVEnv   string `envconfig:"summary_v_env"`
	GnapiDomain   string `envconfig:"gnapi_domain"`

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
