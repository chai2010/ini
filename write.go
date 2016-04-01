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
)

func (c *Config) Save(fname string, header string) error {
	file, err := os.Create(fname)
	if err != nil {
		return err
	}

	buf := bufio.NewWriter(file)
	if err = c.write(buf, header); err != nil {
		return err
	}
	buf.Flush()

	return file.Close()
}

func (c *Config) WriteTo(w io.Writer, header string) error {
	buf := bufio.NewWriter(w)
	if err := c.write(buf, header); err != nil {
		return err
	}
	if err := buf.Flush(); err != nil {
		return err
	}
	return nil
}

// WriteFile saves the configuration representation to a file.
// The desired file permissions must be passed as in os.Open. The header is a
// string that is saved as a comment in the first line of the file.
func (c *Config) WriteFile(fname string, perm os.FileMode, header string) error {
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	buf := bufio.NewWriter(file)
	if err = c.write(buf, header); err != nil {
		return err
	}
	buf.Flush()

	return file.Close()
}

func (c *Config) write(buf *bufio.Writer, header string) (err error) {
	if header != "" {
		header = strings.Replace(header, "\r\n", "\n", -1)
		header = strings.Replace(header, "\n", "\r\n", -1)

		// Add comment character after of each new line.
		if i := strings.Index(header, "\n"); i != -1 {
			header = strings.Replace(header, "\n", "\n"+c.comment, -1)
		}

		if _, err = buf.WriteString(c.comment + header + "\n"); err != nil {
			return err
		}
	}

	for _, section := range c.GetSectionList() {
		options := c.GetSectionKeyList(section)

		// Skip default section if empty.
		if section == DEFAULT_SECTION && len(options) == 0 {
			continue
		}

		if _, err = buf.WriteString("\r\n[" + section + "]\r\n"); err != nil {
			return err
		}

		for _, option := range options {
			s := fmt.Sprint(option, c.separator, c.dataMap[section][option].v, "\r\n")
			if _, err = buf.WriteString(s); err != nil {
				return err
			}
		}
	}

	if _, err = buf.WriteString("\r\n"); err != nil {
		return err
	}
	return nil
}
