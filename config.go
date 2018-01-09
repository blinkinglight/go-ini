package ini

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type item struct {
	haveValue bool
	value     string
	comment   string
}

const keyNotFound = "key not found in config file"

type Config struct {
	keys   []string
	config map[string]map[string]item
	mu     sync.RWMutex
}

func New() *Config {
	c := new(Config)
	c.config = make(map[string]map[string]item)
	c.keys = []string{}
	return c
}

func spl(q string) []string {
	r := []string{}
	r = strings.Split(q, ".")
	return r
}

func (c *Config) Read(filename string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	f, err := os.Open(filename)

	if err != nil {
		panic("failed open config file")
	}

	c.keys = []string{}
	c.config = make(map[string]map[string]item)

	rs := bufio.NewScanner(f)

	openB := false

	tag := ""
	lineN := 0
	for rs.Scan() {
		lineN++
		t := strings.Trim(rs.Text(), " \n")

		if len(t) > 0 {
			if t[0] == '#' {
				continue
			}
		}

		if len(t) > 0 && !(strings.Contains(t, "{") || strings.Contains(t, "}")) && !openB {
			panic("Syntax error on line: " + strconv.Itoa(lineN))

		}

		if strings.Contains(t, "{") && openB {
			panic("Syntax error on line: " + strconv.Itoa(lineN))
		}

		if strings.Contains(t, "{") {
			tag = t[0:strings.Index(t, "{")]
			tag = strings.Trim(tag, " ")
			openB = true
			c.config[tag] = make(map[string]item)
			c.keys = append(c.keys, tag)
			continue
		}

		if strings.Contains(t, "}") && !openB {
			panic("Syntax error on line: " + strconv.Itoa(lineN))
		}

		if strings.Contains(t, "}") {
			openB = false
			tag = ""
			continue
		}

		if tag != "" {
			item := item{}

			comment := ""

			if strings.Contains(t, " //") {
				c_tmp := strings.SplitN(t, " //", 2)
				comment = strings.Trim(c_tmp[1], " ")
				t = c_tmp[0]
			}

			if strings.Contains(t, "=") {
				opts := strings.SplitN(t, "=", 2)
				opts[0] = strings.Trim(opts[0], " ")
				opts[1] = strings.Trim(opts[1], " ")
				item.value = opts[1]
				item.haveValue = true
				item.comment = comment
				c.config[tag][opts[0]] = item
			} else {
				item.value = ""
				item.haveValue = false
				item.comment = comment
				c.config[tag][t] = item
			}
		}

	}
}

func (c *Config) GetKeysDefault(section string, d []string) []string {
	r, e := c.GetKeys(section)
	if e != nil {
		return d
	}
	return r
}

func (c *Config) GetKeysListDefault(section string, d []string) []string {
	r, e := c.GetKeysList(section)
	if e != nil {
		return d
	}
	return r
}

func (c *Config) GetStringDefault(section, key, d string) string {
	r, e := c.GetString(section, key)
	if e != nil {
		return d
	}
	return r
}

func (c *Config) GetBoolDefault(section, key string, d bool) bool {
	r, e := c.GetBool(section, key)
	if e != nil {
		return d
	}
	return r
}

func (c *Config) GetFloatDefault(section, key string, d float32) float32 {
	r, e := c.GetFloat(section, key)
	if e != nil {
		return d
	}
	return r
}

func (c *Config) GetIntDefault(section, key string, d int) int {
	r, e := c.GetInt(section, key)
	if e != nil {
		return d
	}
	return r
}

func (c *Config) GetFloat64Default(section, key string, d float64) float64 {
	r, e := c.GetFloat64(section, key)
	if e != nil {
		return d
	}
	return r
}

func (c *Config) GetInt64Default(section, key string, d int64) int64 {
	r, e := c.GetInt64(section, key)
	if e != nil {
		return d
	}
	return r
}

func (c *Config) Sections() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	list := []string{}

	for k, _ := range c.config {
		list = append(list, k)
	}
	return list
}

func (c *Config) GetKeys(section string) ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	list := []string{}
	if _, ok := c.config[section]; !ok {
		return []string{}, fmt.Errorf(keyNotFound)
	}
	for k, _ := range c.config[section] {
		list = append(list, k)
	}
	return list, nil
}

func (c *Config) GetKeysList(section string) ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	list := []string{}

	if _, ok := c.config[section]; !ok {
		return []string{}, fmt.Errorf(keyNotFound)
	}

	for k, v := range c.config[section] {
		if v.haveValue == true {
			continue
		}
		list = append(list, k)
	}

	return list, nil
}

func (c *Config) GetString(section, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	r, ok := c.config[section][key]
	if !ok {
		return "", fmt.Errorf(keyNotFound)
	}
	return r.value, nil
}

func (c *Config) GetBool(section, key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.config[section][key]
	if !ok {
		return false, fmt.Errorf(keyNotFound)
	}
	r, e := strconv.ParseBool(v.value)
	if e != nil {
		return false, fmt.Errorf(keyNotFound)
	}
	return r, nil
}

func (c *Config) GetInt(section, key string) (int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	tmp, ok := c.config[section][key]
	if !ok {
		return 0, fmt.Errorf(keyNotFound)
	}

	r, e := strconv.ParseInt(tmp.value, 16, 32)

	if e != nil {
		return 0, fmt.Errorf("%v", e)
	}
	return int(r), nil
}

func (c *Config) GetFloat(section, key string) (float32, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	tmp, ok := c.config[section][key]
	if !ok {
		return 0, fmt.Errorf(keyNotFound)
	}

	r, e := strconv.ParseFloat(tmp.value, 64)

	if e != nil {
		return 0, e
	}
	return float32(r), nil
}

func (c *Config) GetInt64(section, key string) (int64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	tmp, ok := c.config[section][key]
	if !ok {
		return 0, fmt.Errorf(keyNotFound)
	}

	r, e := strconv.ParseInt(tmp.value, 16, 64)

	if e != nil {
		return 0, fmt.Errorf("%v", e)
	}
	return r, nil
}

func (c *Config) GetFloat64(section, key string) (float64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	tmp, ok := c.config[section][key]
	if !ok {
		return 0, fmt.Errorf(keyNotFound)
	}

	r, e := strconv.ParseFloat(tmp.value, 64)

	if e != nil {
		return 0, e
	}
	return r, nil
}

func (c *Config) Delete(section, key string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(key) == 0 {
		delete(c.config, section)
		ks := []string{}
		for _, v := range c.keys {
			if v == section {
				continue
			}
			ks = append(ks, v)
		}
		c.keys = ks
	} else {
		delete(c.config[section], key)
	}
}

func (c *Config) Set(section, key string, _value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	i := item{}
	value := fmt.Sprintf("%v", _value)

	if len(value) == 0 {
		i.value = ""
		i.haveValue = false
	} else {
		i.value = value
		i.haveValue = true
	}

	if _, ok := c.config[section]; !ok {
		c.config[section] = make(map[string]item)
		c.keys = append(c.keys, section)
	}

	c.config[section][key] = i

}

func (c *Config) Exists(section, key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, ok := c.config[section]; ok {
		if _, ok := c.config[section][key]; ok {
			return true
		}
	}
	return false
}

func (c *Config) Write(filename string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	f, err := os.Create(filename)
	if err != nil {
		panic("failed save config")
	}

	for _, section := range c.keys {
		f.WriteString(fmt.Sprintln(section, "{"))
		for k, v := range c.config[section] {
			if !v.haveValue {
				continue
			}
			comment := ""
			if len(v.comment) > 0 {
				comment = "// " + v.comment
			}
			f.WriteString(fmt.Sprintf("  %s = %s %s\n", k, v.value, comment))
		}

		for k, v := range c.config[section] {
			if v.haveValue {
				continue
			}
			comment := ""
			if len(v.comment) > 0 {
				comment = "// " + v.comment
			}
			f.WriteString(fmt.Sprintf("  %s %s\n", k, comment))
		}
		f.WriteString(fmt.Sprintln("}"))
		f.WriteString(fmt.Sprintln())
	}

}
