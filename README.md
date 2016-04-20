# Gannett Newsfetch

Cache news articles from Gannett APIs and cache them in a Mongo store for easy serving.

## Setup

1. Install Go dependencies:

        go get

2. Build the Go binary:

        go build

3. Set up environment variables. It is recommended that you copy `.env.sample` into `.env` and adjust as necessary. Apply them via `source .env` or, better yet, use [autoenv](https://github.com/horosgrisa/autoenv).

4. Set up Python environment for the summarizer:

        which virtualenv || pip install virtualenv
        virtualenv summary_venv
        pip install -r requirements.txt
        (cd summary_venv/bin/; ln -s ../../summary.py)


## Environment variables

* `MONGO_URI` - DB uri for mongo store
* `SITE_CODES` - a comma-separated list of Gannett Site codes
* `GANNETT_API_KEY` - API key for Gannett API
* `SUMMARY_V_ENV` - absolute path to the virutal environment for summarization


## Commands
### Articles
Fetch gannett news articles for the current day and summarizes
```
$ ./gannett-newsfetch articles
```

