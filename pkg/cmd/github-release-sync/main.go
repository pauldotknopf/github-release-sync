package main

import (
	"fmt"
	"os"

	"github.com/urface/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "github-release-sync"
	app.Usage = "A tool that syncs a directory with a github release."
	app.Version = Version
	app.HideVersion = true

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "github-release-sync: %s\n", err)
		os.Exit(1)
	}
}
