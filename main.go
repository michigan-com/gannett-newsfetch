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
	"github.com/michigan-com/brvty-api/mongoqueue"
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
			if config.MongoURI != "" {
				var err error
				session, err = SetupMongoSession(config.MongoURI)
				if err != nil {
					log.Fatalf("Failed to connect to '%s': %v", config.MongoURI, err)
					os.Exit(ExitCodeErrDependencies)
				}
				defer session.Close()
			}

			commands.GetArticles(session, config.SiteCodes, config.GannettAPIKey)
		},
	}
	root.AddCommand(articlesCmd)

	scrapeCommand := &cobra.Command{
		Use:   "scrape-and-summarize",
		Short: "Grab stories that we see in chartbeat but not the Gannett API",
		Run: func(command *cobra.Command, args []string) {
			config.LoopInterval = time.Duration(loopSec) * time.Second

			if config.MongoURI == "" {
				log.Fatalf("No mongo uri specified")
				os.Exit(ExitCodeErrConfig)
			}

			if config.GannettAssetAPIKey == "" {
				log.Warning("Gannett API Key needs to be set (env.gannett_asset_api_key)")
			}

			session, err := SetupMongoSession(config.MongoURI)
			if err != nil {
				log.Fatalf("Failed to connect to '%s': %v", config.MongoURI, err)
				os.Exit(ExitCodeErrDependencies)
			}
			defer session.Close()

			var client *brvtyclient.Client
			if config.BrvtyURL != "" {
				client = brvtyclient.New(config.BrvtyURL, config.BrvtyAPIKey)
			}

			queue, err := newQueue(session)
			if err != nil {
				log.Fatalf("Failed to initialize queue: %v", err)
				os.Exit(ExitCodeErrDependencies)
			}

			commands.ScrapeAndSummarize(session, client, queue, config.BrvtyTimeout, config.LoopInterval, config.MongoURI, config.SummaryVEnv, config.GannettAssetAPIKey)
		},
	}
	root.AddCommand(scrapeCommand)

	breakingCommand := &cobra.Command{
		Use:   "breaking-news",
		Short: "Check the Gannett API for breaking news",
		Run: func(command *cobra.Command, argv []string) {
			if config.MongoURI == "" {
				log.Fatalf("No mongo uri specified")
				os.Exit(ExitCodeErrConfig)
			}
			if len(config.SiteCodes) == 0 {
				log.Fatalf("No site codes input, please set the SITE_CODES env variable")
				os.Exit(ExitCodeErrConfig)
			}
			config.LoopInterval = time.Duration(loopSec) * time.Second

			session, err := SetupMongoSession(config.MongoURI)
			if err != nil {
				log.Fatalf("Failed to connect to '%s': %v", config.MongoURI, err)
				os.Exit(ExitCodeErrDependencies)
			}
			defer session.Close()

			commands.FetchBreakingNews(session, config.SiteCodes, config.GNAPIDomain, config.LoopInterval, config.GannettAPIKey)
		},
	}
	root.AddCommand(breakingCommand)

	queueCommand := &cobra.Command{
		Use:   "run-jobs",
		Short: "Run pending queued jobs",
		Run: func(command *cobra.Command, argv []string) {
			if config.MongoURI == "" {
				log.Fatalf("No mongo uri specified")
				os.Exit(ExitCodeErrConfig)
			}

			if config.BrvtyURL == "" {
				log.Fatalf("No Brvty URL specified")
				os.Exit(ExitCodeErrConfig)
			}

			session, err := SetupMongoSession(config.MongoURI)
			if err != nil {
				log.Fatalf("Failed to connect to '%s': %v", config.MongoURI, err)
				os.Exit(ExitCodeErrDependencies)
			}
			defer session.Close()

			client := brvtyclient.New(config.BrvtyURL, config.BrvtyAPIKey)

			queue, err := newQueue(session)
			if err != nil {
				log.Fatalf("Failed to initialize queue: %v", err)
				os.Exit(ExitCodeErrDependencies)
			}

			commands.RunQueuedJobs(session, client, queue, config.BrvtyTimeout)
		},
	}
	root.AddCommand(queueCommand)

	root.PersistentFlags().IntVarP(&loopSec, "loop", "l", -1, "Time in seconds to sleep before looping and hitting the apis again")

	log.Infof(`Running Gannett Newsfetch for Site Codes: %v`, config.SiteCodes)

	root.Execute()
}

func newQueue(session *mgo.Session) (*mongoqueue.Queue, error) {
	queue := mongoqueue.New(session.DB("").C("queue"), mongoqueue.Params{
		Logger: log.New(),
	})

	err := queue.Migrate()
	if err != nil {
		return nil, err
	}

	return queue, err
}
