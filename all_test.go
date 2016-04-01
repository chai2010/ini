// Copyright 2016 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ini

import (
	"bufio"
	"os"
	"reflect"
	"strings"
	"testing"
)

const (
	tmpFilename    = "testdata/__test.go"
	sourceFilename = "testdata/source.ini"
	targetFilename = "testdata/target.ini"
)

func testGet(t *testing.T, c *Config, section string, option string,
	expected interface{}) {
	ok := false
	switch expected.(type) {
	case string:
		v, _ := c.GetString(section, option)
		if v == expected.(string) {
			ok = true
		}
	case int:
		v, _ := c.GetInt(section, option)
		if v == expected.(int) {
			ok = true
		}
	case bool:
		v, _ := c.GetBool(section, option)
		if v == expected.(bool) {
			ok = true
		}
	default:
		t.Fatalf("Bad test case")
	}
	if !ok {
		v, _ := c.GetString(section, option)
		t.Errorf("Get failure: expected different value for %s %s (expected: [%#v] got: [%#v])", section, option, expected, v)
	}
}

// TestInMemory creates configuration representation and run multiple tests in-memory.
func TestInMemory(t *testing.T) {
	c := New(nil)

	// == Test empty structure

	// should be empty
	if len(c.GetSectionList()) != 1 {
		t.Errorf("Sections failure: invalid length")
	}

	// test presence of missing section
	if c.HasSection("no-section") {
		t.Errorf("HasSection failure: invalid section")
	}

	// get options for missing section
	if c.HasSection("no-section") {
		t.Errorf("Options failure: invalid section")
	}

	// test presence of option for missing section
	if c.HasSectionKey("no-section", "no-option") {
		t.Errorf("HasSection failure: invalid/section/option")
	}

	// get value from missing section/option
	_, err := c.GetString("no-section", "no-option")
	if err == nil {
		t.Errorf("String failure: got value for missing section/option")
	}

	// get value from missing section/option
	_, err = c.GetInt("no-section", "no-option")
	if err == nil {
		t.Errorf("Int failure: got value for missing section/option")
	}

	// remove missing section
	if c.RemoveSection("no-section") {
		t.Errorf("RemoveSection failure: removed missing section")
	}

	// remove missing section/option
	if c.RemoveSectionKey("no-section", "no-option") {
		t.Errorf("RemoveOption failure: removed missing section/option")
	}

	// == Fill up structure

	// add section
	if !c.AddSection("section1") {
		t.Errorf("AddSection failure: false on first insert")
	}

	// re-add same section
	if c.AddSection("section1") {
		t.Errorf("AddSection failure: true on second insert")
	}

	// default section always exists
	if c.AddSection(DEFAULT_SECTION) {
		t.Errorf("AddSection failure: true on default section insert")
	}

	// add option/value
	if !c.AddSectionKey("section1", "option1", "value1") {
		t.Errorf("AddOption failure: false on first insert")
	}
	testGet(t, c, "section1", "option1", "value1") // read it back

	// overwrite value
	if c.AddSectionKey("section1", "option1", "value2") {
		t.Errorf("AddOption failure: true on second insert")
	}
	testGet(t, c, "section1", "option1", "value2") // read it back again

	// remove option/value
	if !c.RemoveSectionKey("section1", "option1") {
		t.Errorf("RemoveOption failure: false on first remove")
	}

	// remove again
	if c.RemoveSectionKey("section1", "option1") {
		t.Errorf("RemoveOption failure: true on second remove")
	}

	// read it back again
	_, err = c.GetString("section1", "option1")
	if err == nil {
		t.Errorf("String failure: got value for removed section/option")
	}

	// remove existing section
	if !c.RemoveSection("section1") {
		t.Errorf("RemoveSection failure: false on first remove")
	}

	// remove again
	if c.RemoveSection("section1") {
		t.Errorf("RemoveSection failure: true on second remove")
	}

	// == Test types

	// add section
	if !c.AddSection("section2") {
		t.Errorf("AddSection failure: false on first insert")
	}

	// add number
	if !c.AddSectionKey("section2", "test-number", "666") {
		t.Errorf("AddOption failure: false on first insert")
	}
	testGet(t, c, "section2", "test-number", 666) // read it back

	// add 'yes' (bool)
	if !c.AddSectionKey("section2", "test-yes", "yes") {
		t.Errorf("AddOption failure: false on first insert")
	}
	testGet(t, c, "section2", "test-yes", true) // read it back

	// add 'false' (bool)
	if !c.AddSectionKey("section2", "test-false", "false") {
		t.Errorf("AddOption failure: false on first insert")
	}
	testGet(t, c, "section2", "test-false", false) // read it back

	// == Test cycle

	c.AddSectionKey(DEFAULT_SECTION, "opt1", "%(opt2)s")
	c.AddSectionKey(DEFAULT_SECTION, "opt2", "%(opt1)s")

	_, err = c.GetString(DEFAULT_SECTION, "opt1")
	if err == nil {
		t.Errorf("String failure: no error for cycle")
	} else if strings.Index(err.Error(), "cycle") < 0 {
		t.Errorf("String failure: incorrect error for cycle")
	}
}

// TestReadFile creates a 'tough' configuration file and test (read) parsing.
func TestReadFile(t *testing.T) {
	file, err := os.Create(tmpFilename)
	if err != nil {
		t.Fatal("Test cannot run because cannot write temporary file: " + tmpFilename)
	}

	err = os.Setenv("GO_CONFIGFILE_TEST_ENV_VAR", "configvalue12345")
	if err != nil {
		t.Fatalf("Test cannot run because cannot set environment variable GO_CONFIGFILE_TEST_ENV_VAR: %#v", err)
	}

	buf := bufio.NewWriter(file)
	buf.WriteString("optionInDefaultSection=true\n")
	buf.WriteString("[section-1]\n")
	buf.WriteString("option1=value1 ; This is a comment\n")
	buf.WriteString("option2 : 2#Not a comment\t#Now this is a comment after a TAB\n")
	buf.WriteString("  # Let me put another comment\n")
	buf.WriteString("option3= line1\n line2: \n\tline3=v # Comment multiline with := in value\n")
	buf.WriteString("; Another comment\n")
	buf.WriteString("[" + DEFAULT_SECTION + "]\n")
	buf.WriteString("variable1=small\n")
	buf.WriteString("variable2=a_part_of_a_%(variable1)s_test\n")
	buf.WriteString("[secTION-2]\n")
	buf.WriteString("IS-flag-TRUE=Yes\n")
	buf.WriteString("[section-1] # comment on section header\n") // continue again [section-1]
	buf.WriteString("option4=this_is_%(variable2)s.\n")
	buf.WriteString("envoption1=this_uses_${GO_CONFIGFILE_TEST_ENV_VAR}_env\n")
	buf.WriteString("optionInDefaultSection=false")
	buf.Flush()
	file.Close()

	c, err := Load(tmpFilename, nil)
	if err != nil {
		t.Fatalf("ReadDefault failure: %s", err)
	}

	// check number of sections
	if len(c.GetSectionList()) != 3 {
		t.Errorf("Sections failure: wrong number of sections")
	}

	// check number of options 6 of [section-1] plus 2 of [default]
	opts := c.GetSectionKeyList("section-1")
	if len(opts) != 6 {
		t.Errorf("Options failure: wrong number of options: %d", len(opts))
	}

	testGet(t, c, "section-1", "option1", "value1")
	testGet(t, c, "section-1", "option2", "2#Not a comment")
	testGet(t, c, "section-1", "option3", "line1\nline2:\nline3=v")
	testGet(t, c, "section-1", "option4", "this_is_a_part_of_a_small_test.")
	testGet(t, c, "section-1", "envoption1", "this_uses_configvalue12345_env")
	testGet(t, c, "section-1", "optionInDefaultSection", false)
	testGet(t, c, "section-2", "optionInDefaultSection", true)
	testGet(t, c, "secTION-2", "IS-flag-TRUE", true) // case-sensitive
}

// TestWriteReadFile tests writing and reading back a configuration file.
func TestWriteReadFile(t *testing.T) {
	cw := New(nil)

	// write file; will test only read later on
	cw.AddSection("First-Section")
	cw.AddSectionKey("First-Section", "option1", "value option1")
	cw.AddSectionKey("First-Section", "option2", "2")

	cw.AddSectionKey("", "host", "www.example.com")
	cw.AddSectionKey(DEFAULT_SECTION, "protocol", "https://")
	cw.AddSectionKey(DEFAULT_SECTION, "base-url", "%(protocol)s%(host)s")

	cw.AddSectionKey("Another-Section", "useHTTPS", "y")
	cw.AddSectionKey("Another-Section", "url", "%(base-url)s/some/path")

	cw.WriteFile(tmpFilename, 0644, "Test file for test-case")

	// read back file and test
	cr, err := Load(tmpFilename, nil)
	if err != nil {
		t.Fatalf("ReadDefault failure: %s", err)
	}

	testGet(t, cr, "First-Section", "option1", "value option1")
	testGet(t, cr, "First-Section", "option2", 2)
	testGet(t, cr, "Another-Section", "useHTTPS", true)
	testGet(t, cr, "Another-Section", "url", "https://www.example.com/some/path")

	defer os.Remove(tmpFilename)
}

// TestSectionOptions tests read options in a section without default options.
func TestSectionOptions(t *testing.T) {
	cw := New(nil)

	// write file; will test only read later on
	cw.AddSection("First-Section")
	cw.AddSectionKey("First-Section", "option1", "value option1")
	cw.AddSectionKey("First-Section", "option2", "2")

	cw.AddSectionKey("", "host", "www.example.com")
	cw.AddSectionKey(DEFAULT_SECTION, "protocol", "https://")
	cw.AddSectionKey(DEFAULT_SECTION, "base-url", "%(protocol)s%(host)s")

	cw.AddSectionKey("Another-Section", "useHTTPS", "y")
	cw.AddSectionKey("Another-Section", "url", "%(base-url)s/some/path")

	cw.WriteFile(tmpFilename, 0644, "Test file for test-case")

	// read back file and test
	cr, err := Load(tmpFilename, nil)
	if err != nil {
		t.Fatalf("ReadDefault failure: %s", err)
	}

	options := cr.GetSectionKeyList("First-Section")

	if len(options) != 2 {
		t.Fatalf("SectionOptions reads wrong data: %v", options)
	}

	expected := map[string]bool{
		"option1": true,
		"option2": true,
	}
	actual := map[string]bool{}

	for _, v := range options {
		actual[v] = true
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("SectionOptions reads wrong data: %v", options)
	}

	options = cr.GetSectionKeyList(DEFAULT_SECTION)

	expected = map[string]bool{
		"host":     true,
		"protocol": true,
		"base-url": true,
	}
	actual = map[string]bool{}

	for _, v := range options {
		actual[v] = true
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("SectionOptions reads wrong data: %v", options)
	}

	defer os.Remove(tmpFilename)
}

// TestMerge tests merging 2 configurations.
func TestMerge(t *testing.T) {
	target, error := Load(targetFilename, nil)
	if error != nil {
		t.Fatalf("Unable to read target config file '%s'", targetFilename)
	}

	source, error := Load(sourceFilename, nil)
	if error != nil {
		t.Fatalf("Unable to read source config file '%s'", sourceFilename)
	}

	target.MergeFrom(source)

	// Assert whether a regular option was merged from source -> target
	if result, _ := target.GetString(DEFAULT_SECTION, "one"); result != "source1" {
		t.Errorf("Expected 'one' to be '1' but instead it was '%s'", result)
	}
	// Assert that a non-existent option in source was not overwritten
	if result, _ := target.GetString(DEFAULT_SECTION, "five"); result != "5" {
		t.Errorf("Expected 'five' to be '5' but instead it was '%s'", result)
	}
	// Assert that a folded option was correctly unfolded
	if result, _ := target.GetString(DEFAULT_SECTION, "two_+_three"); result != "source2 + source3" {
		t.Errorf("Expected 'two_+_three' to be 'source2 + source3' but instead it was '%s'", result)
	}
	if result, _ := target.GetString(DEFAULT_SECTION, "four"); result != "4" {
		t.Errorf("Expected 'four' to be '4' but instead it was '%s'", result)
	}

	// Assert that a section option has been merged
	if result, _ := target.GetString("X", "x.one"); result != "sourcex1" {
		t.Errorf("Expected '[X] x.one' to be 'sourcex1' but instead it was '%s'", result)
	}
	if result, _ := target.GetString("X", "x.four"); result != "x4" {
		t.Errorf("Expected '[X] x.four' to be 'x4' but instead it was '%s'", result)
	}
}
