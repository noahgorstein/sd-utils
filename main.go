package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

func main() {
	sdUtils := &StardogUtils{}
	ctx := kong.Parse(sdUtils,
		kong.Name("sd-utils"),
		kong.Description("A collection of miscellaneous utilities for Stardog."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: false,
		}))

	if err := ctx.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
