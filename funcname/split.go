package funcname

import "strings"

// Split splits s into two strings by the first '.' after the last '/'.
// In case of finding no '/', Split splits s by the first '.'.
// In case of finding no '.', Split returns pkgname = s and name = "".
func Split(s string) (pkgname, name string) {
	i := strings.LastIndexByte(s, '/') + 1
	j := strings.IndexByte(s[i:], '.') + 1
	if j == 0 {
		return s, ""
	}
	k := i + j
	return s[:k-1], s[k:]
}
