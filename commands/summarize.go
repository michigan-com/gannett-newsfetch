package commands

import (
	"encoding/json"
	"fmt"
	"os/exec"

	log "github.com/Sirupsen/logrus"

	"github.com/michigan-com/gannett-newsfetch/config"
	"github.com/michigan-com/gannett-newsfetch/lib"
)

type SummaryResponse struct {
	Skipped    int `json:"skipped"`
	Summarized int `json:"summarized"`
}

/*
	Run a python process to summarize all articles in the ToSummarize collection
*/
func ProcessSummaries(toSummarize []interface{}, mongoUri string) (*SummaryResponse, error) {
	summResp := &SummaryResponse{}

	session := lib.DBConnect(mongoUri)
	bulk := session.DB("").C("ToSummarize").Bulk()
	bulk.Upsert(toSummarize...)
	_, err := bulk.Run()
	if err != nil {
		return summResp, err
	}

	envConfig, _ := config.GetEnv()

	if envConfig.SummaryVEnv == "" {
		return nil, fmt.Errorf("Missing SUMMARY_VENV environtment variable, skipping summarizer")
	}

	cmd := fmt.Sprintf("%s/bin/python", envConfig.SummaryVEnv)
	pyScript := fmt.Sprintf("%s/bin/summary.py", envConfig.SummaryVEnv)

	log.Infof("Executing command: %s %s %s", cmd, pyScript, mongoUri)

	out, err := exec.Command(cmd, pyScript, mongoUri).Output()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(out, summResp); err != nil {
		return nil, err
	}
	fmt.Println(summResp)

	return summResp, nil
}
