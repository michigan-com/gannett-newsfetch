package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/michigan-com/brvty-api/mongoqueue"
	"github.com/michigan-com/gannett-newsfetch/commands"
)

func main() {
	var brvtyTimeout time.Duration
	var interval time.Duration
	var verboseMgo bool
	flag.DurationVar(&interval, "l", 0, "Run continuously with this interval")
	flag.DurationVar(&brvtyTimeout, "brvty-timeout", 30*time.Second, "Set Brvty API timeout")
	flag.BoolVar(&verboseMgo, "verbose-mgo", false, "Enable verbose logging in mgo")
	flag.Parse()

	log.SetOutput(os.Stderr)
	log.SetFlags(log.Ltime)
	mgo.SetLogger(log.New(os.Stderr, "[mgo] ", log.Ltime))
	queueLogger := log.New(os.Stderr, "", log.Ltime)
	mgo.SetDebug(verboseMgo)

	runtime.GOMAXPROCS(runtime.NumCPU())

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		fmt.Fprintf(os.Stderr, "** missing MONGO_URI env variable\n")
		os.Exit(64) // EX_USAGE
	}

	// brvtyURL := os.Getenv("BRVTY_URL")
	// if brvtyURL == "" {
	// 	fmt.Fprintf(os.Stderr, "** missing BRVTY_URL env variable\n")
	// 	os.Exit(64) // EX_USAGE
	// }
	// brvtyAPIKey := os.Getenv("BRVTY_API_KEY")
	// if brvtyAPIKey == "" {
	// 	fmt.Fprintf(os.Stderr, "** missing BRVTY_API_KEY env variable\n")
	// 	os.Exit(64) // EX_USAGE
	// }
	//
	// summaryVEnv := os.Getenv("SUMMARY_V_ENV")
	// if summaryVEnv == "" {
	// 	fmt.Fprintf(os.Stderr, "** missing SUMMARY_V_ENV env variable\n")
	// 	os.Exit(64) // EX_USAGE
	// }

	assetAPIKey := os.Getenv("GANNETT_ASSET_API_KEY")
	if assetAPIKey == "" {
		fmt.Fprintf(os.Stderr, "** missing GANNETT_ASSET_API_KEY env variable\n")
		os.Exit(64) // EX_USAGE
	}

	session, err := mgo.Dial(mongoURI)
	if err != nil {
		log.Printf("Failed to connect to %#v: %v", mongoURI, err)
		os.Exit(1)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	// client := brvtyclient.New(brvtyURL, brvtyAPIKey)
	queue := mongoqueue.New(session.DB("").C("queue"), mongoqueue.Params{
		Logger: queueLogger,
	})

	commands.ScrapeAndSummarize(session, queue, interval, mongoURI, assetAPIKey)
}
