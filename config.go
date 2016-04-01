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

type Options struct {
	Comment   string // default is DEFAULT_COMMENT
	Separator string // default is ALTERNATIVE_SEPARATOR
	PreSpace  bool   // default is true
	PostSpace bool   // default is true
}

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
//	opt.Comment: has to be `DEFAULT_COMMENT` or `ALTERNATIVE_COMMENT`
//	opt.Separator: has to be `DEFAULT_SEPARATOR` or `ALTERNATIVE_SEPARATOR`
//	opt.PreSpace: indicate if is inserted a space before of the separator
//	opt.PostSpace: indicate if is added a space after of the separator
//
func New(opt *Options) *Config {
	if opt == nil {
		opt = &Options{
			Comment:   DEFAULT_COMMENT,
			Separator: DEFAULT_SEPARATOR,
			PreSpace:  true,
			PostSpace: true,
		}
	}
	if opt.Comment == "" {
		opt.Comment = DEFAULT_COMMENT
	}
	if opt.Separator == "" {
		opt.Separator = DEFAULT_SEPARATOR
	}

	if opt.Comment != DEFAULT_COMMENT && opt.Comment != ALTERNATIVE_COMMENT {
		panic("invalid comment:" + opt.Comment)
	}
	if opt.Separator != DEFAULT_SEPARATOR && opt.Separator != ALTERNATIVE_SEPARATOR {
		panic("invalid separator:" + opt.Separator)
	}

	comment := opt.Comment
	separator := opt.Separator

	// Get spaces around separator
	if opt.PreSpace {
		separator = " " + separator
	}
	if opt.PostSpace {
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

// Merge merges the given configuration "source" with this one ("p").
//
// Merging means that any option (under any section) from source that is not in
// p will be copied into p. When the p already has an option with
// the same name and section then it is overwritten (i.o.w. the source wins).
//
func (p *Config) Merge(source *Config) {
	if source == nil || len(source.dataMap) == 0 {
		return
	}

	for section, option := range source.dataMap {
		for optionName, optionValue := range option {
			p.AddSectionKey(section, optionName, optionValue.v)
		}
	}
}
