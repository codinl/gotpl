package tpl

import (
	"bytes"
)

func Child() []byte {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\n    child\n\n\naaa")
	return _buffer.Bytes()
}
