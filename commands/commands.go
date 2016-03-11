package commands

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/michigan-com/gannett-newsfetch/config"
)

var (
	VERSION    string
	COMMITHASH string
)

var NewsfetchCommand = &cobra.Command{
	Use: "newsfetch",
}

func Run(version, commit string) {
	VERSION = version
	COMMITHASH = commit
	AddCommands()
	PrepareEnvironment()
	NewsfetchCommand.Execute()
}

/*
	Add all necessary command line commands
*/
func AddCommands() {
	NewsfetchCommand.AddCommand(articlesCmd)
}

/*
	Prepare the environemtn for newsfetch. Read in the env variables, doing some
	basic env var checking to make sure they're set.
*/
func PrepareEnvironment() {
	env, _ := config.GetEnv()

	siteCodeSplit := strings.Split(env.SiteCodes, ",")
	log.Info(env)
	siteCodes := make([]string, 0, len(siteCodeSplit))
	for _, code := range siteCodeSplit {
		if code != "" {
			siteCodes = append(siteCodes, code)
		}
	}

	if len(siteCodes) == 0 {
		log.Fatal("No site codes input, please set the SITE_CODES env variable")
	}

	log.Info(fmt.Sprintf(`

	Running Gannett Newsfetch

		Site Codes: %v

	`, siteCodeSplit))
}
