// Copyright 2016 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ini

import (
	"regexp"
)

const (
	// Default section name.
	DEFAULT_SECTION   = "DEFAULT"
	DEFAULT_COMMENT   = "# "
	DEFAULT_SEPARATOR = "="

	ALTERNATIVE_COMMENT   = "; "
	ALTERNATIVE_SEPARATOR = ":"

	// Maximum allowed depth when recursively substituing variable names.
	_DEPTH_VALUES = 64
)

var (
	// Strings accepted as boolean.
	boolString = map[string]bool{
		"t":     true,
		"true":  true,
		"y":     true,
		"yes":   true,
		"on":    true,
		"1":     true,
		"f":     false,
		"false": false,
		"n":     false,
		"no":    false,
		"off":   false,
		"0":     false,

		"是": true,
		"否": false,
		"不": false,
	}

	varRegExp    = regexp.MustCompile(`%\(([a-zA-Z0-9_.\-]+)\)s`) // %(variable)s
	envVarRegExp = regexp.MustCompile(`\${([a-zA-Z0-9_.\-]+)}`)   // ${envvar}
)

// Config is the representation of configuration settings.
type Config struct {
	comment   string
	separator string

	// Sections order
	lastIdSection   int            // Last section identifier
	idSectionMap    map[string]int // Section : position
	lastIdOptionMap map[string]int // Section : last identifier

	// Section -> option : value
	dataMap map[string]map[string]*tValue
}

// tValue holds the input position for a value.
type tValue struct {
	position int    // Option order
	v        string // value
}

// New creates an empty configuration representation.
// This representation can be filled with AddSection and AddOption and then
// saved to a file using WriteFile.
//
//	comment: has to be `DEFAULT_COMMENT` or `ALTERNATIVE_COMMENT`
//	separator: has to be `DEFAULT_SEPARATOR` or `ALTERNATIVE_SEPARATOR`
//	preSpace: indicate if is inserted a space before of the separator
//	postSpace: indicate if is added a space after of the separator
//
func New(comment, separator string, preSpace, postSpace bool) *Config {
	if comment != DEFAULT_COMMENT && comment != ALTERNATIVE_COMMENT {
		panic("invalid comment:" + comment)
	}

	if separator != DEFAULT_SEPARATOR && separator != ALTERNATIVE_SEPARATOR {
		panic("invalid separator:" + separator)
	}

	// Get spaces around separator
	if preSpace {
		separator = " " + separator
	}
	if postSpace {
		separator += " "
	}

	c := new(Config)

	c.comment = comment
	c.separator = separator
	c.idSectionMap = make(map[string]int)
	c.lastIdOptionMap = make(map[string]int)
	c.dataMap = make(map[string]map[string]*tValue)

	c.AddSection(DEFAULT_SECTION) // Default section always exists.

	return c
}

// NewDefault creates a configuration representation with values by default.
func NewDefault() *Config {
	return New(DEFAULT_COMMENT, DEFAULT_SEPARATOR, true, true)
}

// Merge merges the given configuration "source" with this one ("target").
//
// Merging means that any option (under any section) from source that is not in
// target will be copied into target. When the target already has an option with
// the same name and section then it is overwritten (i.o.w. the source wins).
//
func (target *Config) Merge(source *Config) {
	if source == nil || len(source.dataMap) == 0 {
		return
	}

	for section, option := range source.dataMap {
		for optionName, optionValue := range option {
			target.AddOption(section, optionName, optionValue.v)
		}
	}
}
