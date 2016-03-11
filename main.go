package main

import (
	"runtime"

	"github.com/michigan-com/gannett-newsfetch/commands"
)

// Version number that gets compiled via `make build` or `make install`
var VERSION string

// Git commit hash that gets compiled via `make build` or `make install`
var COMMITHASH string

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	commands.Run(VERSION, COMMITHASH)
}
