package wxbot

import "regexp"

// FindString 数据匹配
func FindString(pattern, value, key string) string {
	Map := SelectString(pattern, value)
	if val, ok := Map[key]; ok {
		return val
	}
	return ""
}

// SelectString 数组匹配
func SelectString(pattern, value string) map[string]string {
	Exp := myRegexp{regexp.MustCompile(pattern)}
	Map := Exp.FindStringSubmatchMap(value)
	return Map
}

type myRegexp struct {
	*regexp.Regexp
}

func (r *myRegexp) FindStringSubmatchMap(s string) map[string]string {
	captures := make(map[string]string)

	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}

	for i, name := range r.SubexpNames() {
		if i == 0 {
			continue
		}
		captures[name] = match[i]
	}
	return captures
}
