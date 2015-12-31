package tpl

import (
	"bytes"
)

func Base() []byte {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\naaa")
	return _buffer.Bytes()
}
