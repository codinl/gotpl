package tpl

import (
	"bytes"
)

func Base() []byte {
	var _buffer bytes.Buffer
	_buffer.WriteString("@{{\n}}\n@block  basebase  {{}}\n\naaaaaa")
	return _buffer.Bytes()
}
