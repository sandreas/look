package actions

import (
	"bufio"
	"errors"
	"github.com/sandreas/log"
	"github.com/urfave/cli"
	"io"
	"os"
	"regexp"
)

type At struct {
}

func (action *At) Execute(c *cli.Context) error {
	settings := parseActionParams(c)
	initLogging(settings)
	defer log.Flush()

	len := c.Args().Len()
	lastArg := ""
	var reader *bufio.Reader
	if len > 0 {
		lastArg = c.Args().Slice()[len-1]
		input, err := os.OpenFile(lastArg, os.O_RDONLY, 0755)

		if err != nil {
			log.Warn("Could not open file", lastArg)
			return err
		}
		reader = bufio.NewReader(input)
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			reader = bufio.NewReader(os.Stdin)
		} else {
			log.Warn("Either a file or stdin has to be provided")
			return errors.New("Either a file or stdin has to be provided")
		}
	}

	// file reader
	// os.Stdin reader

	var matchers []*regexp.Regexp
	filters := c.StringSlice("filter")
	for _, f := range filters {
		re1, err := regexp.Compile(f)
		if err != nil {
			log.Warn("Invalid regex", f)
			return err
		}

		matchers = append(matchers, re1)
	}

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Warn("Could not read line")
			return err
		}

		skipLine := false
		for _, m := range matchers {
			if !m.Match(line) {
				skipLine = true
				break
			}
		}
		if skipLine {
			continue
		}
		println(string(line))
	}

	return nil
}
