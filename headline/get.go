// Package headline はテキストの先頭の行を返す関数を提供する.
package headline

import "strings"

// Get は空白ではない最初の行を返す.
// 改行は '\n'
// 返される文字列に改行を含まない.
func Get(s string) string {
	for {
		if s == "" {
			return s
		}
		i := strings.IndexByte(s, '\n')
		if i == -1 {
			return s
		}
		if v := s[:i]; strings.TrimSpace(v) != "" {
			return v
		}
		s = s[i+1:]
	}
}
