package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2"

	"github.com/michigan-com/brvty-api/brvtyclient"
	"github.com/michigan-com/gannett-newsfetch/commands"
)

// Version number that gets compiled via `make build` or `make install`
var VERSION string

// Git commit hash that gets compiled via `make build` or `make install`
var COMMITHASH string

const (
	ExitCodeErrProcessing   = 1
	ExitCodeErrDependencies = 2
	ExitCodeErrConfig       = 3
)

func main() {
	var loopSec int

	runtime.GOMAXPROCS(runtime.NumCPU())

	config, err := ParseConfig()
	if err != nil {
		log.Fatalf(`Error loading environment config: %v`, err)
		os.Exit(ExitCodeErrConfig)
	}
	fmt.Printf("%#v\n", config)

	root := &cobra.Command{
		Use: "newsfetch",
	}

	articlesCmd := &cobra.Command{
		Use:   "articles",
		Short: "Get articles published on the current day from the Gannett API",
		Run: func(command *cobra.Command, args []string) {
			if len(config.SiteCodes) == 0 {
				log.Fatal("Please set the SITE_CODES env variable")
				os.Exit(ExitCodeErrConfig)
			}

			var session *mgo.Session
			if config.MongoUri != "" {
				var err error
				session, err = SetupMongoSession(config.MongoUri)
				if err != nil {
					log.Fatalf("Failed to connect to '%s': %v", config.MongoUri, err)
					os.Exit(ExitCodeErrDependencies)
				}
			}
			defer closeSessionIfNotNil(session)

			commands.GetArticles(session, config.SiteCodes, config.GannettApiKey)
		},
	}
	root.AddCommand(articlesCmd)

	scrapeCommand := &cobra.Command{
		Use:   "scrape-and-summarize",
		Short: "Grab stories that we see in chartbeat but not the Gannett API",
		Run: func(command *cobra.Command, args []string) {
			config.LoopInterval = time.Duration(loopSec) * time.Second

			session, err := SetupMongoSession(config.MongoUri)
			if err != nil {
				log.Fatalf("Failed to connect to '%s': %v", config.MongoUri, err)
				os.Exit(ExitCodeErrDependencies)
			}
			defer session.Close()

			var client *brvtyclient.Client
			if config.BrvtyURL != "" {
				client = brvtyclient.New(config.BrvtyURL, config.BrvtyAPIKey)
			}

			commands.ScrapeAndSummarize(session, client, config.BrvtyTimeout, config.LoopInterval, config.MongoUri, config.SummaryVEnv)
		},
	}
	root.AddCommand(scrapeCommand)

	breakingCommand := &cobra.Command{
		Use:   "breaking-news",
		Short: "Check the Gannett API for breaking news",
		Run: func(command *cobra.Command, argv []string) {
			if config.MongoUri == "" {
				log.Fatalf("No mongo uri specified")
				os.Exit(ExitCodeErrConfig)
			}
			if len(config.SiteCodes) == 0 {
				log.Fatalf("No site codes input, please set the SITE_CODES env variable")
				os.Exit(ExitCodeErrConfig)
			}
			config.LoopInterval = time.Duration(loopSec) * time.Second

			session, err := SetupMongoSession(config.MongoUri)
			if err != nil {
				log.Fatalf("Failed to connect to '%s': %v", config.MongoUri, err)
				os.Exit(ExitCodeErrDependencies)
			}
			defer session.Close()

			commands.FetchBreakingNews(session, config.SiteCodes, config.GnapiDomain, config.LoopInterval, config.GannettApiKey)
		},
	}
	root.AddCommand(breakingCommand)

	root.PersistentFlags().IntVarP(&loopSec, "loop", "l", -1, "Time in seconds to sleep before looping and hitting the apis again")

	log.Infof(`Running Gannett Newsfetch for Site Codes: %v`, config.SiteCodes)

	root.Execute()
}
