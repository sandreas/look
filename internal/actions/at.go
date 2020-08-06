package actions

import (
	"bufio"
	"errors"
	"github.com/sandreas/log"
	"github.com/urfave/cli"
	"io"
	"os"
	"regexp"
	"time"
)

const FlagFilter = "expression"
const FlagWatch = "watch"
const FlagLineCount = "lines"

type At struct {
}

func (action *At) Execute(c *cli.Context) error {
	settings := parseActionParams(c)
	initLogging(settings)
	defer log.Flush()

	argsLen := c.Args().Len()
	lastArg := ""
	var reader *bufio.Reader
	if argsLen > 0 {
		lastArg = c.Args().Slice()[argsLen-1]
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
			log.Warn("either a file or stdin has to be provided")
			return errors.New("either a file or stdin has to be provided")
		}
	}

	// file reader
	// os.Stdin reader

	var matchers []*regexp.Regexp
	filters := c.StringSlice(FlagFilter)
	for _, f := range filters {
		re1, err := regexp.Compile(f)
		if err != nil {
			log.Warn("Invalid regex", f)
			return err
		}

		matchers = append(matchers, re1)
	}

	maxLines := c.Int(FlagLineCount)
	if maxLines < 0 {
		maxLines = 0
	}

	var lineBuffer = make([][]byte, maxLines)
	shouldWatchFile := c.Bool(FlagWatch)
OuterLoop:
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				if maxLines != 0 {
					flushLineBuffer(lineBuffer)
					lineBuffer = lineBuffer[:0]
				}
				if shouldWatchFile {
					time.Sleep(2 * time.Second)
					continue OuterLoop
				}
				break
			}
			log.Warn("Could not read line")
			return err
		}

		for _, m := range matchers {
			if !m.Match(line) {
				continue OuterLoop
			}
		}

		if maxLines == 0 {
			flushLine(line)
		} else {
			length := len(lineBuffer)
			if length == maxLines {
				lineBuffer = lineBuffer[1:]
			}
			lineBuffer = append(lineBuffer, line)
		}
	}

	return nil
}

func flushLineBuffer(lineBuffer [][]byte) {
	for _, line := range lineBuffer {
		flushLine(line)
	}
}

func flushLine(line []byte) {
	println(string(line))
}
