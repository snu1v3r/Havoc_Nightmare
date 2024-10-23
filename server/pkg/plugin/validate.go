package plugin

import "regexp"

func ValidName(name string) bool {
	return regexp.MustCompile("^[a-zA-Z0-9_.-]+$").MatchString(name)
}

func ValidType(_type string) bool {
	switch _type {
	case TypeListener:
		fallthrough
	case TypeAgent:
		fallthrough
	case TypeManagement:
		return true
	default:
		return false
	}
}
