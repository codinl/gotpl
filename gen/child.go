package tpl

import (
	"bytes"
)

func Child() []byte {
	var _buffer bytes.Buffer
	_buffer.WriteString("@{{\n}}\n@block  basebase  {{\n        childchild\n}}")
	return _buffer.Bytes()
}
