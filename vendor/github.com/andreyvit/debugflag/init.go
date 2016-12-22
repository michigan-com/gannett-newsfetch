package debugflag

import (
	"os"
)

const EnvVar = "DEBUG"

var enabled map[string]bool

func init() {
	spec := os.Getenv(EnvVar)
	Reset(spec)
}

func Reset(spec string) {
	e, d := parse(spec)

	enabled = make(map[string]bool)

	for _, el := range e {
		enabled[el] = true
	}
	for _, el := range d {
		enabled[el] = false
	}
}
