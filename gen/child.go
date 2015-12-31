package tpl

import (
	"bytes"
	"github.com/codinl/gotpl/template"
)

func Child(val int) []byte {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\n<html>\n<body>\n\n    <div>base aa content</div>\n\n    \n    <div>child bb block content.</div>\n\n\n    \n    ")
	for i := 0; i < 10; i++ {

		_buffer.WriteString("<p>")
		_buffer.WriteString(gotpl.HTMLEscape(i))
		_buffer.WriteString("</p>\n    ")
	}
	_buffer.WriteString("\n\n\n\n<div>this is TestSection content. Param \"val\" is: ")
	_buffer.WriteString(gotpl.HTMLEscape(val))
	_buffer.WriteString(" </div>\n\n\n</body>\n\n</html>")
	return _buffer.Bytes()
}
