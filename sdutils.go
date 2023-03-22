package main

import (
	"github.com/noahgorstein/sd-utils/concurrent-select"
)

// StardogUtils is the command-line interface for sd-utils
type StardogUtils struct {
	BenchSelect concurrentselect.Options `cmd:"" help:"execute a SPARQL select query concurrently to benchmark server response times"`
}
