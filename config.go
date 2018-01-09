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

var keys []string
var config map[string]map[string]item = make(map[string]map[string]item)

var mu sync.RWMutex

func spl(q string) []string {
	r := []string{}
	r = strings.Split(q, ".")
	return r
}

func Read(filename string) {
	mu.Lock()
	defer mu.Unlock()
	f, err := os.Open(filename)

	if err != nil {
		panic("failed open config file")
	}

	keys = []string{}
	config = make(map[string]map[string]item)

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
			config[tag] = make(map[string]item)
			keys = append(keys, tag)
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
				config[tag][opts[0]] = item
			} else {
				item.value = ""
				item.haveValue = false
				item.comment = comment
				config[tag][t] = item
			}
		}

	}
}

func GetKeysDefault(section string, d []string) []string {
	r, e := GetKeys(section)
	if e != nil {
		return d
	}
	return r
}

func GetKeysListDefault(section string, d []string) []string {
	r, e := GetKeysList(section)
	if e != nil {
		return d
	}
	return r
}

func GetStringDefault(section, key, d string) string {
	r, e := GetString(section, key)
	if e != nil {
		return d
	}
	return r
}

func GetBoolDefault(section, key string, d bool) bool {
	r, e := GetBool(section, key)
	if e != nil {
		return d
	}
	return r
}

func GetFloatDefault(section, key string, d float32) float32 {
	r, e := GetFloat(section, key)
	if e != nil {
		return d
	}
	return r
}

func GetIntDefault(section, key string, d int) int {
	r, e := GetInt(section, key)
	if e != nil {
		return d
	}
	return r
}

func GetFloat64Default(section, key string, d float64) float64 {
	r, e := GetFloat64(section, key)
	if e != nil {
		return d
	}
	return r
}

func GetInt64Default(section, key string, d int64) int64 {
	r, e := GetInt64(section, key)
	if e != nil {
		return d
	}
	return r
}

func Sections() []string {
	mu.RLock()
	defer mu.RUnlock()
	list := []string{}

	for k, _ := range config {
		list = append(list, k)
	}
	return list
}

func GetKeys(section string) ([]string, error) {
	mu.RLock()
	defer mu.RUnlock()
	list := []string{}
	if _, ok := config[section]; !ok {
		return []string{}, fmt.Errorf(keyNotFound)
	}
	for k, _ := range config[section] {
		list = append(list, k)
	}
	return list, nil
}

func GetKeysList(section string) ([]string, error) {
	mu.RLock()
	defer mu.RUnlock()
	list := []string{}

	if _, ok := config[section]; !ok {
		return []string{}, fmt.Errorf(keyNotFound)
	}

	for k, v := range config[section] {
		if v.haveValue == true {
			continue
		}
		list = append(list, k)
	}

	return list, nil
}

func GetString(section, key string) (string, error) {
	mu.RLock()
	defer mu.RUnlock()
	r, ok := config[section][key]
	if !ok {
		return "", fmt.Errorf(keyNotFound)
	}
	return r.value, nil
}

func GetBool(section, key string) (bool, error) {
	mu.RLock()
	defer mu.RUnlock()
	v, ok := config[section][key]
	if !ok {
		return false, fmt.Errorf(keyNotFound)
	}
	r, e := strconv.ParseBool(v.value)
	if e != nil {
		return false, fmt.Errorf(keyNotFound)
	}
	return r, nil
}

func GetInt(section, key string) (int, error) {
	mu.RLock()
	defer mu.RUnlock()
	tmp, ok := config[section][key]
	if !ok {
		return 0, fmt.Errorf(keyNotFound)
	}

	r, e := strconv.ParseInt(tmp.value, 16, 32)

	if e != nil {
		return 0, fmt.Errorf("%v", e)
	}
	return int(r), nil
}

func GetFloat(section, key string) (float32, error) {
	mu.RLock()
	defer mu.RUnlock()
	tmp, ok := config[section][key]
	if !ok {
		return 0, fmt.Errorf(keyNotFound)
	}

	r, e := strconv.ParseFloat(tmp.value, 64)

	if e != nil {
		return 0, e
	}
	return float32(r), nil
}

func GetInt64(section, key string) (int64, error) {
	mu.RLock()
	defer mu.RUnlock()
	tmp, ok := config[section][key]
	if !ok {
		return 0, fmt.Errorf(keyNotFound)
	}

	r, e := strconv.ParseInt(tmp.value, 16, 64)

	if e != nil {
		return 0, fmt.Errorf("%v", e)
	}
	return r, nil
}

func GetFloat64(section, key string) (float64, error) {
	mu.RLock()
	defer mu.RUnlock()
	tmp, ok := config[section][key]
	if !ok {
		return 0, fmt.Errorf(keyNotFound)
	}

	r, e := strconv.ParseFloat(tmp.value, 64)

	if e != nil {
		return 0, e
	}
	return r, nil
}

func Delete(section, key string) {
	mu.RLock()
	defer mu.RUnlock()
	if len(key) == 0 {
		delete(config, section)
		ks := []string{}
		for _, v := range keys {
			if v == section {
				continue
			}
			ks = append(ks, v)
		}
		keys = ks
	} else {
		delete(config[section], key)
	}
}

func Set(section, key string, _value interface{}) {
	mu.Lock()
	defer mu.Unlock()

	i := item{}
	value := fmt.Sprintf("%v", _value)

	if len(value) == 0 {
		i.value = ""
		i.haveValue = false
	} else {
		i.value = value
		i.haveValue = true
	}

	if _, ok := config[section]; !ok {
		config[section] = make(map[string]item)
		keys = append(keys, section)
	}

	config[section][key] = i

}

func Exists(section, key string) bool {
	mu.RLock()
	defer mu.RUnlock()
	if _, ok := config[section]; ok {
		if _, ok := config[section][key]; ok {
			return true
		}
	}
	return false
}

func Write(filename string) {
	mu.Lock()
	defer mu.Unlock()

	f, err := os.Create(filename)
	if err != nil {
		panic("failed save config")
	}

	for _, section := range keys {
		f.WriteString(fmt.Sprintln(section, "{"))
		for k, v := range config[section] {
			if !v.haveValue {
				continue
			}
			comment := ""
			if len(v.comment) > 0 {
				comment = "// " + v.comment
			}
			f.WriteString(fmt.Sprintf("  %s = %s %s\n", k, v.value, comment))
		}

		for k, v := range config[section] {
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
