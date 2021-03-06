// Copyright 2016 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ini

func (c *Config) MustValue(section, key string, defaultVal ...string) string {
	v, err := c.GetValue(section, key)
	if err != nil {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		} else {
			panic(err)
		}
	}
	return v
}

func (c *Config) MustBool(section, key string, defaultVal ...bool) bool {
	v, err := c.GetBool(section, key)
	if err != nil {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		} else {
			panic(err)
		}
	}
	return v
}

func (c *Config) MustInt(section, key string, defaultVal ...int) int {
	v, err := c.GetInt(section, key)
	if err != nil {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		} else {
			panic(err)
		}
	}
	return v
}

func (c *Config) MustFloat64(section, key string, defaultVal ...float64) float64 {
	v, err := c.GetFloat64(section, key)
	if err != nil {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		} else {
			panic(err)
		}
	}
	return v
}

func (c *Config) MustString(section, key string, defaultVal ...string) string {
	v, err := c.GetString(section, key)
	if err != nil {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		} else {
			panic(err)
		}
	}
	return v
}
