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

func Load(fname string, opt *Options) (c *Config, err error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	br := bufio.NewReader(file)
	if bom, _ := br.Peek(3); string(bom) == "\xEF\xBB\xBF" {
		br.Discard(3)
	}

	c = New(opt)
	if err = c.read(br); err != nil {
		return nil, err
	}

	if err = file.Close(); err != nil {
		return nil, err
	}

	return c, nil
}

func LoadFrom(r io.Reader, opt *Options) (c *Config, err error) {
	var br *bufio.Reader
	if p, ok := r.(*bufio.Reader); ok {
		br = p
	} else {
		br = bufio.NewReader(r)
	}

	if bom, _ := br.Peek(3); string(bom) == "\xEF\xBB\xBF" {
		br.Discard(3)
	}

	c = New(opt)
	if err = c.read(br); err != nil {
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
			prev, _ := c.GetValue(section, option)
			value := strings.TrimSpace(l)
			c.AddSectionKey(section, option, prev+"\n"+value)

		// Other alternatives
		default:
			i := strings.IndexAny(l, "=:")

			switch {
			// Option and value
			case i > 0 && l[0] != ' ' && l[0] != '\t': // found an =: and it's not a multiline continuation
				option = strings.TrimSpace(l[0:i])
				value := strings.TrimSpace(l[i+1:])
				c.AddSectionKey(section, option, value)

			default:
				return fmt.Errorf("ini: could not parse line: %d", l)
			}
		}
	}
	return scanner.Err()
}
