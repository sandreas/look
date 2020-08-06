package main

import (
	"fmt"
	"github.com/sandreas/look/internal/actions"
	"github.com/urfave/cli"
	"os"
)

func main() {
	globalFlags := []cli.Flag{
		&cli.BoolFlag{Name: "quiet, q", Usage: "do not show any output"},                                           // does quiet make sense in find?
		&cli.BoolFlag{Name: "force, f", Usage: "force the requested action - even if it might be not a good idea"}, // does force make sense in find?
		&cli.BoolFlag{Name: "debug", Usage: "debug mode with logging to Stdout and into $HOME/.graft/application.log"},
	}

	atFlags := []cli.Flag{
		// &cli.BoolFlag{Name: "keep-duplicates", Usage: "keep duplicate files"},
		// &cli.StringFlag{Name: "tpl", Usage: "filename template"},
		// &cli.StringFlag{Name: "include-media-types", Value: "image,video", Usage: "media types to include"},

		&cli.StringSliceFlag{
			Name:    actions.FlagExpression, // todo rename to pattern?
			Aliases: []string{"e"},
		},
		&cli.StringSliceFlag{
			Name:    actions.FlagReplacements,
			Aliases: []string{"r"},
		},
		&cli.BoolFlag{
			Name:    actions.FlagWatch,
			Aliases: []string{"w"},
		},

		&cli.IntFlag{
			Name:    actions.FlagLineCount,
			Aliases: []string{"l"},
		},
	}

	app := cli.NewApp()
	app.Name = "look"
	app.Version = "0.1"
	app.Usage = "look at log files"

	app.Commands = []*cli.Command{
		{
			Name:    "at",
			Aliases: []string{"a"},
			Action:  new(actions.At).Execute,
			Usage:   "look at file",
			Flags:   mergeFlags(globalFlags, atFlags),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		println(fmt.Errorf("error: %s", err))
	}
}

func mergeFlags(flagsToMerge ...[]cli.Flag) []cli.Flag {
	var mergedFlags []cli.Flag
	for _, flags := range flagsToMerge {
		mergedFlags = append(mergedFlags, flags...)
	}
	return mergedFlags
}
