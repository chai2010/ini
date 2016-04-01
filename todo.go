// Copyright 2016 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ini

import (
	"io"
)

func Load(filename string, more ...string) (c *Config, err error) {
	panic("TODO")
}

func LoadFromData(data []byte) (c *Config, err error) {
	panic("TODO")
}

func LoadFromReader(in io.Reader) (c *Config, err error) {
	panic("TODO")
}

func (c *Config) DeleteKey(section, key string) bool {
	panic("TODO")
}
func (c *Config) DeleteSection(section string) bool {
	panic("TODO")
}
func (c *Config) Float64(section, key string) (float64, error) {
	panic("TODO")
}
func (c *Config) GetKeyComments(section, key string) (comments string) {
	panic("TODO")
}
func (c *Config) GetKeyList(section string) []string {
	panic("TODO")
}
func (c *Config) GetSection(section string) (map[string]string, error) {
	panic("TODO")
}
func (c *Config) GetSectionComments(section string) (comments string) {
	panic("TODO")
}
func (c *Config) GetSectionList() []string {
	panic("TODO")
}
func (c *Config) GetValue(section, key string) (string, error) {
	panic("TODO")
}
func (c *Config) Int64(section, key string) (int64, error) {
	panic("TODO")
}
func (c *Config) MustBool(section, key string, defaultVal ...bool) bool {
	panic("TODO")
}
func (c *Config) MustFloat64(section, key string, defaultVal ...float64) float64 {
	panic("TODO")
}
func (c *Config) MustInt(section, key string, defaultVal ...int) int {
	panic("TODO")
}
func (c *Config) MustInt64(section, key string, defaultVal ...int64) int64 {
	panic("TODO")
}
func (c *Config) MustValue(section, key string, defaultVal ...string) string {
	panic("TODO")
}
func (c *Config) MustValueArray(section, key, delim string) []string {
	panic("TODO")
}
func (c *Config) MustValueRange(section, key, defaultVal string, candidates []string) string {
	panic("TODO")
}
func (c *Config) MustValueSet(section, key string, defaultVal ...string) (string, bool) {
	panic("TODO")
}
func (c *Config) SetKeyComments(section, key, comments string) bool {
	panic("TODO")
}
func (c *Config) SetSectionComments(section, comments string) bool {
	panic("TODO")
}
