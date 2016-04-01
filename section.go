// Copyright 2016 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ini

// HasSection checks if the configuration has the given section.
// (The default section always exists.)
func (c *Config) HasSection(section string) bool {
	if section == "" || section == DEFAULT_SECTION {
		return true
	}
	_, ok := c.dataMap[section]
	return ok
}

// AddSection adds a new section to the configuration.
//
// If the section is nil then uses the section by default which it's already
// created.
//
// It returns true if the new section was inserted, and false if the section
// already existed.
func (c *Config) AddSection(section string) bool {
	if section == "" {
		section = DEFAULT_SECTION
	}
	if _, ok := c.dataMap[section]; ok {
		return false
	}

	c.dataMap[section] = make(map[string]*tValue)

	// Section order
	c.idSectionMap[section] = c.lastIdSection
	c.lastIdSection++

	return true
}

// RemoveSection removes a section from the configuration.
// It returns true if the section was removed, and false if section did not exist.
func (c *Config) RemoveSection(section string) bool {
	// Default section cannot be removed.
	if section == "" || section == DEFAULT_SECTION {
		return false
	}
	if _, ok := c.dataMap[section]; !ok {
		return false
	}

	delete(c.dataMap, section)
	delete(c.lastIdOptionMap, section)
	delete(c.idSectionMap, section)
	return true
}

// GetSectionList returns the list of sections in the configuration.
// (The default section always exists).
func (c *Config) GetSectionList() (sections []string) {
	for i := 0; i < c.lastIdSection; i++ {
		for section, id := range c.idSectionMap {
			if id == i {
				sections = append(sections, section)
			}
		}
	}
	return
}
