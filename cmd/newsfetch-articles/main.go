package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/michigan-com/gannett"
	"github.com/michigan-com/gannett-newsfetch/commands"
)

func main() {
	var interval time.Duration
	var domains gannett.DomainList
	var verboseMgo bool
	flag.DurationVar(&interval, "l", 0, "Run continuously with this interval")
	flag.BoolVar(&verboseMgo, "verbose-mgo", false, "Enable verbose logging in mgo")
	flag.Var(&domains, "d", "Only enable these domains (comma-separated and/or multiple options)")
	flag.Parse()

	log.SetOutput(os.Stderr)
	log.SetFlags(log.Ltime)
	mgo.SetLogger(log.New(os.Stderr, "[mgo] ", log.Ltime))
	mgo.SetDebug(verboseMgo)

	runtime.GOMAXPROCS(runtime.NumCPU())

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		fmt.Fprintf(os.Stderr, "** missing MONGO_URI env variable\n")
		os.Exit(64) // EX_USAGE
	}

	gannettAPIKey := os.Getenv("GANNETT_SEARCH_API_KEY")
	if gannettAPIKey == "" {
		fmt.Fprintf(os.Stderr, "** missing GANNETT_SEARCH_API_KEY env variable\n")
		os.Exit(64) // EX_USAGE
	}

	if len(domains) == 0 {
		domains = gannett.AllDomains()
	}

	var siteCodes []string
	for _, d := range domains {
		pub := gannett.PublicationByDomain(d)
		siteCodes = append(siteCodes, pub.SiteCode)
	}

	session, err := mgo.Dial(mongoURI)
	if err != nil {
		log.Printf("Failed to connect to %#v: %v", mongoURI, err)
		os.Exit(1)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	log.Printf("✓ Enabled %d domains", len(domains))

	for {
		commands.GetArticles(session, siteCodes, gannettAPIKey)
		log.Printf("✓ Done")
		if interval == 0 {
			break
		}
		log.Printf("Sleeping for %v ms...", interval.Nanoseconds()/1000000)
		time.Sleep(interval)
	}
}
