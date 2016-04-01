// Copyright 2016 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ini

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type Options struct {
	Comment   string
	Separator string
	PreSpace  bool
	PostSpace bool
}

func Load(fname string, opt *Options) (c *Config, err error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	if opt != nil {
		c = New(opt.Comment, opt.Separator, opt.PreSpace, opt.PostSpace)
	} else {
		c = NewDefault()
	}

	if err = c.read(bufio.NewReader(file)); err != nil {
		return nil, err
	}

	if err = file.Close(); err != nil {
		return nil, err
	}

	return c, nil
}

func LoadFrom(r io.Reader, opt *Options) (c *Config, err error) {
	if opt != nil {
		c = New(opt.Comment, opt.Separator, opt.PreSpace, opt.PostSpace)
	} else {
		c = NewDefault()
	}

	if err = c.read(bufio.NewReader(r)); err != nil {
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
