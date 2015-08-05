package tpl

import (
	"bytes"
	"github.com/codinl/gotpl/gotpl"
)

func Test() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("<html>\n\n\n    aaaa\n\n    \n     extends bbb\n\n\n    \n     ")
	for i := 0; i < 10; i++ {

		_buffer.WriteString("<p>")
		_buffer.WriteString(gotpl.HTMLEscape(i))
		_buffer.WriteString("</p>")

	}
	_buffer.WriteString("\n\n\n\ncurPage int    <div>curPage is: ")
	_buffer.WriteString(gotpl.HTMLEscape(curPage))
	_buffer.WriteString(" </div>\n\n\n</html>")

	return _buffer.String()
}
