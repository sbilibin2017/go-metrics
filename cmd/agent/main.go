package main

import (
	"go-metrics/cmd/agent/app"
	"go-metrics/pkg/cli"
	"os"
)

func main() {
	cmd := app.NewCommand()
	code := cli.Run(cmd)
	os.Exit(code)
}
