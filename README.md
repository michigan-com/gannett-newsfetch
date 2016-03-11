# Gannett Newsfetch
Cache news articles from Gannett APIs and cache them in a Mongo store for easy serving

## Setup
```
$ go get
$ go build
```

### Env vars
*See [`test_env.sh`](https://github.com/michigan-com/gannett-newsfetch/blob/master/test_env.sh) for required env variables*

* `MONGO_URI` - DB uri for mongo store
* `SITE_CODES` - a comma-separated list of Gannett Site codes
* `GANNETT_API_KEY` - API key for Gannett API

## Commands
### Articles
Fetch gannett news articles
```
$ ./gannett-newsfetch articles
```


