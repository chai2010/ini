// Copyright 2016 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ini

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Read reads a configuration file and returns its representation.
// All arguments, except `fname`, are related to `New()`
func Read(fname string, comment, separator string, preSpace, postSpace bool) (*Config, error) {
	return _read(fname, New(comment, separator, preSpace, postSpace))
}

// ReadDefault reads a configuration file and returns its representation.
// It uses values by default.
func ReadDefault(fname string) (*Config, error) {
	return _read(fname, NewDefault())
}

// _read is the base to read a file and get the configuration representation.
// That representation can be queried with GetString, etc.
func _read(fname string, c *Config) (*Config, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	if err = c.read(bufio.NewReader(file)); err != nil {
		return nil, err
	}

	if err = file.Close(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) read(buf *bufio.Reader) (err error) {
	var section, option string
	var scanner = bufio.NewScanner(buf)
	for scanner.Scan() {
		l := strings.TrimRightFunc(stripComments(scanner.Text()), unicode.IsSpace)

		// Switch written for readability (not performance)
		switch {
		// Empty line and comments
		case len(l) == 0, l[0] == '#', l[0] == ';':
			continue

		// New section. The [ must be at the start of the line
		case l[0] == '[' && l[len(l)-1] == ']':
			option = "" // reset multi-line value
			section = strings.TrimSpace(l[1 : len(l)-1])
			c.AddSection(section)

		// Continuation of multi-line value
		// starts with whitespace, we're in a section and working on an option
		case section != "" && option != "" && (l[0] == ' ' || l[0] == '\t'):
			prev, _ := c.RawString(section, option)
			value := strings.TrimSpace(l)
			c.AddOption(section, option, prev+"\n"+value)

		// Other alternatives
		default:
			i := strings.IndexAny(l, "=:")

			switch {
			// Option and value
			case i > 0 && l[0] != ' ' && l[0] != '\t': // found an =: and it's not a multiline continuation
				option = strings.TrimSpace(l[0:i])
				value := strings.TrimSpace(l[i+1:])
				c.AddOption(section, option, value)

			default:
				return fmt.Errorf("ini: could not parse line: %d", l)
			}
		}
	}
	return scanner.Err()
}
