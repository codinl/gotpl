package tmp

import (
	"bytes"
	"github.com/codinl/gotpl/gotpl"
)

func Base() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("<html>\n\ncurPage int    <div>curPage is: ")
	_buffer.WriteString(gotpl.HTMLEscape(curPage))
	_buffer.WriteString(" </div>\n\n\n\n    aaaa\n\n    \n        bbb\n    \n\n    \n        ccc\n    \n\n\n</html>")

	return _buffer.String()
}
