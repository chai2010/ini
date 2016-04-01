// Copyright 2016 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ini

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Substitutes values, calculated by callback, on matching regex
func (c *Config) computeVar(
	beforeValue *string, regx *regexp.Regexp, headsz, tailsz int,
	withVar func(*string) string,
) (*string, error) {
	var i int
	computedVal := beforeValue
	for i = 0; i < _DEPTH_VALUES; i++ { // keep a sane depth
		vr := regx.FindStringSubmatchIndex(*computedVal)
		if len(vr) == 0 {
			break
		}

		varname := (*computedVal)[vr[headsz]:vr[headsz+1]]
		varVal := withVar(&varname)
		if varVal == "" {
			return &varVal, fmt.Errorf("ini: Option not found: %s", varname)
		}

		// substitute by new value and take off leading '%(' and trailing ')s'
		//  %(foo)s => headsz=2, tailsz=2
		//  ${foo}  => headsz=2, tailsz=1
		newVal := (*computedVal)[0:vr[headsz]-headsz] + varVal + (*computedVal)[vr[headsz+1]+tailsz:]
		computedVal = &newVal
	}

	if i >= _DEPTH_VALUES {
		retVal := ""
		return &retVal, fmt.Errorf(
			"ini: Possible cycle while unfolding variables: max depth of %d reached",
			_DEPTH_VALUES,
		)
	}

	return computedVal, nil
}

// Bool has the same behaviour as String but converts the response to bool.
// See "boolString" for string values converted to bool.
func (c *Config) Bool(section string, option string) (value bool, err error) {
	if section == "" {
		section = DEFAULT_SECTION
	}

	sv, err := c.String(section, option)
	if err != nil {
		return false, err
	}

	value, ok := boolString[strings.ToLower(sv)]
	if !ok {
		return false, fmt.Errorf("ini: could not parse bool value: %v", sv)
	}

	return value, nil
}

// Float64 has the same behaviour as String but converts the response to float.
func (c *Config) Float64(section string, option string) (value float64, err error) {
	if section == "" {
		section = DEFAULT_SECTION
	}

	sv, err := c.String(section, option)
	if err == nil {
		value, err = strconv.ParseFloat(sv, 64)
	}

	return value, err
}

// Int has the same behaviour as String but converts the response to int.
func (c *Config) Int(section string, option string) (value int, err error) {
	if section == "" {
		section = DEFAULT_SECTION
	}

	sv, err := c.String(section, option)
	if err == nil {
		value, err = strconv.Atoi(sv)
	}

	return value, err
}

// RawString gets the (raw) string value for the given option in the section.
// The raw string value is not subjected to unfolding, which was illustrated in
// the beginning of this documentation.
//
// It returns an error if either the section or the option do not exist.
func (c *Config) RawString(section string, option string) (value string, err error) {
	if section == "" {
		section = DEFAULT_SECTION
	}

	if _, ok := c.dataMap[section]; ok {
		if tValue, ok := c.dataMap[section][option]; ok {
			return tValue.v, nil
		}
	}
	return c.RawStringDefault(option)
}

// RawStringDefault gets the (raw) string value for the given option from the
// DEFAULT section.
//
// It returns an error if the option does not exist in the DEFAULT section.
func (c *Config) RawStringDefault(option string) (value string, err error) {
	if tValue, ok := c.dataMap[DEFAULT_SECTION][option]; ok {
		return tValue.v, nil
	}
	return "", fmt.Errorf("ini: option '%s' not found", option)
}

// String gets the string value for the given option in the section.
// If the value needs to be unfolded (see e.g. %(host)s example in the beginning
// of this documentation), then String does this unfolding automatically, up to
// _DEPTH_VALUES number of iterations.
//
// It returns an error if either the section or the option do not exist, or the
// unfolding cycled.
func (c *Config) String(section string, option string) (value string, err error) {
	if section == "" {
		section = DEFAULT_SECTION
	}

	value, err = c.RawString(section, option)
	if err != nil {
		return "", err
	}

	// % variables
	computedVal, err := c.computeVar(&value, varRegExp, 2, 2, func(varName *string) string {
		lowerVar := *varName
		// search variable in default section as well as current section
		varVal, _ := c.dataMap[DEFAULT_SECTION][lowerVar]
		if _, ok := c.dataMap[section][lowerVar]; ok {
			varVal = c.dataMap[section][lowerVar]
		}
		return varVal.v
	})
	value = *computedVal

	if err != nil {
		return value, err
	}

	// $ environment variables
	computedVal, err = c.computeVar(&value, envVarRegExp, 2, 1, func(varName *string) string {
		return os.Getenv(*varName)
	})
	value = *computedVal
	return value, err
}
