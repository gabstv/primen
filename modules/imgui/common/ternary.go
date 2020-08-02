package common

func TernaryString(v bool, truev, falsev string) string {
	if v {
		return truev
	}
	return falsev
}
