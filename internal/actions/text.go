package actions

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/sandreas/log"
	"github.com/urfave/cli"
	"io"
	"os"
	"regexp"
	"time"
)

const FlagQueryWhere = "where"
const FlagQueryWhereNot = "where-not"
const FlagPattern = "pattern"
const FlagReplacement = "replacement"
const FlagWatch = "watch"
const FlagLineCount = "lines"

type Text struct {
}

func (action *Text) Execute(c *cli.Context) error {
	settings := parseActionParams(c)
	initLogging(settings)
	defer log.Flush()

	// todo https://godoc.org/github.com/yalp/jsonpath

	maxLines := c.Int(FlagLineCount)
	shouldWatchFile := c.Bool(FlagWatch)

	whereFlags := c.StringSlice(FlagQueryWhere)
	whereNotFlags := c.StringSlice(FlagQueryWhereNot)
	patterns := c.StringSlice(FlagPattern)
	replacements := c.StringSlice(FlagReplacement)

	argsLen := c.Args().Len()
	lastArg := ""
	var reader *bufio.Reader
	var input *os.File
	if argsLen > 0 {
		lastArg = c.Args().First()
		i, err := os.OpenFile(lastArg, os.O_RDONLY, 0755)

		if err != nil {
			log.Warn("Could not open file", lastArg)
			return err
		}
		input = i
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
	replacementsLen := len(replacements)
	if replacementsLen > 0 && replacementsLen != len(patterns) {
		log.Warn("replacements must have the same count as patterns")
		return errors.New("replacements must have the same count as patterns")
	}

	for _, f := range patterns {
		re1, err := regexp.Compile(f)
		if err != nil {
			log.Warn("Invalid regex", f)
			return err
		}

		matchers = append(matchers, re1)
	}

	var whereFilters []*regexp.Regexp
	for _, where := range whereFlags {
		re1, err := regexp.Compile(where)
		if err != nil {
			log.Warn("Invalid regex", where)
			return err
		}

		whereFilters = append(whereFilters, re1)
	}

	var whereNotFilters []*regexp.Regexp
	for _, whereNot := range whereNotFlags {
		re1, err := regexp.Compile(whereNot)
		if err != nil {
			log.Warn("Invalid regex", whereNot)
			return err
		}

		whereNotFilters = append(whereNotFilters, re1)
	}

	if maxLines < 0 {
		maxLines = 0
	}

	// var lineBuffer = make([][]byte, 0)
	var lineBuffer []string
	inputLastSize := int64(0)
OuterLoop:
	for {
		if input != nil {
			stat, err := input.Stat()
			if err == nil {
				size := stat.Size()
				if size < inputLastSize {
					reader.Reset(reader)
				}
				inputLastSize = size
			}
		}

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

		// check whereFilters
		for _, f := range whereFilters {
			if !f.Match(line) {
				continue OuterLoop
			}
		}

		for _, f := range whereNotFilters {
			if f.Match(line) {
				continue OuterLoop
			}
		}

		// handle patterns and replacements
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
			//flushLine(string(line))
			//println("-----------------")
			//flushLineBuffer(lineBuffer)
			//println("===================")
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
	_, _ = os.Stdout.WriteString(fmt.Sprintln(line))
}
