# Gannett Newsfetch

Fetch news articles from Gannett APIs and cache them in a Mongo store for easy serving.

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

See [.env.sample](https://github.com/michigan-com/gannett-newsfetch/blob/master/.env.sample) for examples

* `MONGO_URI` - DB uri for mongo store
* `SITE_CODES` - a comma-separated list of Gannett Site codes
* `GANNETT_SEARCH_API_KEY` - API key for Gannett Search API
* `GANNETT_ASSET_API_KEY` - API key for Gannett Asset API
* `SUMMARY_V_ENV` - absolute path to the virutal environment for summarization
* `GNAPI_DOMAIN` - Domain to update when snapshots are saved


## Debugging

Set `DEBUG` to a comma-separated list of these flags to enable additional behaviors:

* `json:articles`: dump incoming article JSONs to stdout


## Commands
### Articles
Fetch gannett news articles for the current day
Env variables used: `MONGO_URI`, `SITE_CODES`, `GANNETT_SEARCH_API_KEY`
```
$ ./gannett-newsfetch articles
```

### Scrape And Summarize
Scrape and summarize articles identified for scraping in `./gannett-newsfetch articles` command
```
$ ./gannett-newsfetch scrape-and-summarize
```

To run indeterminately, add the `-l` flag, indicating the number of seconds to sleep between every loop
```
$ ./gannett-newsfetch scrape-and-summarize -l 5
```

### Breaking News
Get any breaking news alerts from the Gannett API
```
$ ./gannett-newsfetch breaking-news
```

Similar to the scrape and summarize command, add the `-l` flag to loop
```
$ ./gannett-newsfetch breaking-news -l 5
```

## Testing
```
$ go test ./...
```
