package commands

import (
	"encoding/json"
	"fmt"
	"os/exec"

	log "github.com/Sirupsen/logrus"

	"github.com/michigan-com/gannett-newsfetch/config"
)

type SummaryResponse struct {
	Skipped    int `json:"skipped"`
	Summarized int `json:"summarized"`
}

/*
	Run a python process to summarize all articles in the ToSummarize collection
*/
func ProcessSummaries() (*SummaryResponse, error) {
	envConfig, _ := config.GetEnv()

	if envConfig.SummaryVEnv == "" {
		return nil, fmt.Errorf("Missing SUMMARY_VENV environtment variable, skipping summarizer")
	}

	cmd := fmt.Sprintf("%s/bin/python", envConfig.SummaryVEnv)
	pyScript := fmt.Sprintf("%s/bin/summary.py", envConfig.SummaryVEnv)

	log.Infof("Executing command: %s %s %s", cmd, pyScript, envConfig.MongoUri)

	out, err := exec.Command(cmd, pyScript, envConfig.MongoUri).Output()
	if err != nil {
		log.Errorf("\n\nError summarizing articles: %v\n\n", err)
		return nil, err
	}

	summResp := &SummaryResponse{}
	if err := json.Unmarshal(out, summResp); err != nil {
		return nil, err
	}

	return summResp, nil
}
