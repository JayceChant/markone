package main

import (
	"os"

	"github.com/JayceChant/markone/cmd"
	"github.com/urfave/cli"
)

const (
	appName  = "markone"
	appUsage = `Markdown once, export all formats you need.`
	appVer   = "0.1.00"
)

func main() {
	exc := cli.NewApp()
	exc.Name = appName
	exc.Usage = appUsage
	exc.Version = appVer
	exc.Commands = []cli.Command{
		cmd.New,
	}
	exc.Run(os.Args)
}
