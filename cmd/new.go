package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/urfave/cli"
)

const (
	cmdName       = "new"
	layoutDir     = "layout"
	mdSuffix      = ".md"
	defaultLayout = "article"
	timeFormat    = "2006-01-02 15:04:05"
)

var (
	New = cli.Command{
		Name:        cmdName,
		ShortName:   "n",
		Usage:       "New a source file.",
		Description: `New a markdown souce file with markone extended sytax.`,
		Action:      newFile,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "overwrite, o",
				Usage: `Overwrite duplicate file if exists.`,
			},
		},
		ArgsUsage: "[layout] <title>",
	}
)

func newFile(ctx *cli.Context) error {
	na := ctx.NArg()
	if na == 0 {
		cli.ShowCommandHelpAndExit(ctx, cmdName, 0)
	}

	layout := defaultLayout
	title := ""

	switch na {
	case 1:
		title = ctx.Args().Get(0)
	case 2:
		layout = ctx.Args().Get(0)
		title = ctx.Args().Get(1)
	default:
		layout = ctx.Args().Get(0)
		title = strings.Join(ctx.Args()[1:], " ")
	}

	source := path.Join(layoutDir, layout+mdSuffix)
	target := strings.NewReplacer(" ", "-").Replace(title) + mdSuffix

	b, err := ioutil.ReadFile(source)
	if err != nil {
		return cli.NewExitError(
			fmt.Errorf("layout \"%s\" not found:\n\t%v", layout, err), 1)
	}

	tmpl, err := template.New("layout").Funcs(template.FuncMap{
		"title": func() string { return title },
		"date":  func() string { return time.Now().Format(timeFormat) },
	}).Parse(string(b))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	oflag := os.O_WRONLY | os.O_CREATE
	if ctx.Bool("overwrite") {
		oflag |= os.O_TRUNC
	} else {
		oflag |= os.O_EXCL
	}

	f, err := os.OpenFile(target, oflag, 0666)
	defer f.Close()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if err := tmpl.Execute(f, nil); err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}
