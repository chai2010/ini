// Copyright 2016 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ini

import (
	"fmt"
)

// AddOption adds a new option and value to the configuration.
//
// If the section is nil then uses the section by default; if it does not exist,
// it is created in advance.
//
// It returns true if the option and value were inserted, and false if the value
// was overwritten.
func (c *Config) AddOption(section string, option string, value string) bool {
	if section == "" {
		section = DEFAULT_SECTION
	}

	c.AddSection(section) // Make sure section exists

	_, ok := c.dataMap[section][option]
	c.dataMap[section][option] = &tValue{c.lastIdOptionMap[section], value}
	c.lastIdOptionMap[section]++
	return !ok
}

// RemoveOption removes a option and value from the configuration.
// It returns true if the option and value were removed, and false otherwise,
// including if the section did not exist.
func (c *Config) RemoveOption(section string, option string) bool {
	if section == "" {
		section = DEFAULT_SECTION
	}

	if _, ok := c.dataMap[section]; !ok {
		return false
	}

	_, ok := c.dataMap[section][option]
	delete(c.dataMap[section], option)
	return ok
}

// HasOption checks if the configuration has the given option in the section.
// It returns false if either the option or section do not exist.
func (c *Config) HasOption(section string, option string) bool {
	if section == "" {
		section = DEFAULT_SECTION
	}

	if _, ok := c.dataMap[section]; !ok {
		return false
	}

	_, ok := c.dataMap[section][option]
	return ok
}

// Options returns the list of options available in the given section.
// It returns an error if the section does not exist and an empty list if the
// section is empty. Options within the default section are also included.
func (c *Config) Options(section string) (options []string, err error) {
	if section == "" {
		section = DEFAULT_SECTION
	}

	if _, ok := c.dataMap[section]; !ok {
		return nil, fmt.Errorf("ini: section '%s' not found", section)
	}

	// Keep a map of option names we've seen to deduplicate.
	optionMap := make(map[string]struct{},
		len(c.dataMap[DEFAULT_SECTION])+len(c.dataMap[section]))
	for s, _ := range c.dataMap[DEFAULT_SECTION] {
		optionMap[s] = struct{}{}
	}
	for s, _ := range c.dataMap[section] {
		optionMap[s] = struct{}{}
	}

	// Get the keys.
	i := 0
	options = make([]string, len(optionMap))
	for k, _ := range optionMap {
		options[i] = k
		i++
	}

	return options, nil
}

// GetOptionList returns only the list of options available in the given section.
func (c *Config) GetOptionList(section string) (options []string) {
	if section == "" {
		section = DEFAULT_SECTION
	}
	if _, ok := c.dataMap[section]; !ok {
		return
	}

	for i := 0; i < c.lastIdOptionMap[section]; i++ {
		for s, tValue := range c.dataMap[section] {
			if tValue.position == i {
				options = append(options, s)
				break
			}
		}
	}
	return
}
