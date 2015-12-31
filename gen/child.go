package tpl

import (
	"bytes"
)

func Child() []byte {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\n    child\naaa")
	return _buffer.Bytes()
}
