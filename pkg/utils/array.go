package utils

import "bytes"

func (u *Utils) ImplodeStringArr(glue string, pieces []string) string {
	var buf bytes.Buffer
	l := len(pieces)
	for _, str := range pieces {
		buf.WriteString(str)
		if l--; l > 0 {
			buf.WriteString(glue)
		}
	}
	return buf.String()
}

func (u *Utils) ImplodeMapIntString(glue string, pieces map[int]string) string {
	var buf bytes.Buffer
	l := len(pieces)
	for _, str := range pieces {
		buf.WriteString(str)
		if l--; l > 0 {
			buf.WriteString(glue)
		}
	}
	return buf.String()
}

func (u *Utils) InArray(needle interface{}, haystack interface{}) bool {
	switch key := needle.(type) {
	case string:
		for _, item := range haystack.([]string) {
			if key == item {
				return true
			}
		}
	case int:
		for _, item := range haystack.([]int) {
			if key == item {
				return true
			}
		}
	case int64:
		for _, item := range haystack.([]int64) {
			if key == item {
				return true
			}
		}
	default:
		return false
	}
	return false
}
