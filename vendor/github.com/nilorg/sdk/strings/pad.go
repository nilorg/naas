package strings

import (
	"bytes"
)

// PadLeft 填充左边
func PadLeft(str string, l int, s [1]string) string {
	buf := bytes.NewBuffer(nil)
	for i := l; i > 0; i-- {
		buf.WriteString(s[0])
	}
	buf.WriteString(str)
	return buf.String()
}

// PadRight 填充右边
func PadRight(str string, l int, s [1]string) string {
	buf := bytes.NewBufferString(str)
	for i := 0; i < l; i++ {
		buf.WriteString(s[0])
	}
	return buf.String()
}
