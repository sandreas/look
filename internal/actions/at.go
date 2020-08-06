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

const FlagExpression = "expression"
const FlagWatch = "watch"
const FlagLineCount = "lines"
const FlagReplacements = "replacements"

type At struct {
}

func (action *At) Execute(c *cli.Context) error {
	settings := parseActionParams(c)
	initLogging(settings)
	defer log.Flush()

	// todo https://godoc.org/github.com/yalp/jsonpath

	maxLines := c.Int(FlagLineCount)
	shouldWatchFile := c.Bool(FlagWatch)

	argsLen := c.Args().Len()
	lastArg := ""
	var reader *bufio.Reader
	if argsLen > 0 {
		lastArg = c.Args().First()
		input, err := os.OpenFile(lastArg, os.O_RDONLY, 0755)

		if err != nil {
			log.Warn("Could not open file", lastArg)
			return err
		}
		reader = bufio.NewReader(input)
	} else {
		stat, _ := os.Stdin.Stat()
		shouldWatchFile = false
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			reader = bufio.NewReader(os.Stdin)
		} else {
			log.Warn("either a file or stdin has to be provided")
			return errors.New("either a file or stdin has to be provided")
		}
	}

	var matchers []*regexp.Regexp
	expressions := c.StringSlice(FlagExpression)

	replacements := c.StringSlice(FlagReplacements)
	replacementsLen := len(replacements)
	if replacementsLen > 0 && replacementsLen != len(expressions) {
		log.Warn("replacements must have the same count as expressions")
		return errors.New("replacements must have the same count as expressions")
	}

	for _, f := range expressions {
		re1, err := regexp.Compile(f)
		if err != nil {
			log.Warn("Invalid regex", f)
			return err
		}

		matchers = append(matchers, re1)
	}

	if maxLines < 0 {
		maxLines = 0
	}

	// var lineBuffer = make([][]byte, 0)
	var lineBuffer []string
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

		for i, m := range matchers {
			if !m.Match(line) {
				continue OuterLoop
			}

			if replacementsLen > 0 {
				replacement := replacements[i]
				line = m.ReplaceAll(line, []byte(replacement))
			}
		}

		if maxLines == 0 {
			flushLine(string(line))
		} else {
			length := len(lineBuffer)
			if length == maxLines {
				lineBuffer = lineBuffer[1:]
			}
			lineBuffer = append(lineBuffer, string(line))
			flushLine(string(line))
			println("-----------------")
			flushLineBuffer(lineBuffer)
			println("===================")
		}
	}

	return nil
}

func flushLineBuffer(lineBuffer []string) {
	for _, line := range lineBuffer {
		flushLine(line)
	}
}

func flushLine(line string) {
	println(line)
}
