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
	loop       int
)

var NewsfetchCommand = &cobra.Command{
	Use: "newsfetch",
}

func Run(version, commit string) {
	VERSION = version
	COMMITHASH = commit
	AddCommands()
	AddFlags()
	PrepareEnvironment()
	NewsfetchCommand.Execute()
}

/*
	Add all necessary command line commands
*/
func AddCommands() {
	NewsfetchCommand.AddCommand(articlesCmd)
	NewsfetchCommand.AddCommand(cleanupCommand)
}

/*
	Add all necessary flags
*/
func AddFlags() {
	NewsfetchCommand.PersistentFlags().IntVarP(&loop, "loop", "l", -1, "Time in seconds to sleep before looping and hitting the apis again")
}

/*
	Prepare the environemtn for newsfetch. Read in the env variables, doing some
	basic env var checking to make sure they're set.
*/
func PrepareEnvironment() {
	env, _ := config.GetEnv()

	siteCodeSplit := strings.Split(env.SiteCodes, ",")
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
