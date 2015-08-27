package template

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	UNDEF = iota
	AT
	ASSIGN_OPERATOR
	AT_COLON
	AT_STAR_CLOSE       // *@
	AT_STAR_OPEN        // @*
	BACKSLASH
	BRACE_CLOSE         // }
	BRACE_OPEN          // {
	CONTENT
	EMAIL
	ESCAPED_QUOTE       // \"
	HARD_PAREN_CLOSE    // ]
	HARD_PAREN_OPEN     // [
	HTML_TAG_OPEN
	HTML_TAG_CLOSE
	HTML_TAG_VOID_CLOSE
	IDENTIFIER
	KEYWORD
	LOGICAL
	NEWLINE
	NUMERIC_CONTENT
	OPERATOR
	PAREN_CLOSE        // )
	PAREN_OPEN         // (
	PERIOD             // .
	SINGLE_QUOTE       // '
	DOUBLE_QUOTE       // "
	TEXT_TAG_CLOSE
	TEXT_TAG_OPEN
	COMMENT_TAG_OPEN
	COMMENT_TAG_CLOSE
	WHITESPACE
	AT_BLOCK
	AT_SECTION
)

var TokenTypeNames = [...]string{
	"UNDEF", "AT", "ASSIGN_OPERATOR", "AT_COLON",
	"AT_STAR_CLOSE", "AT_STAR_OPEN", "BACKSLASH",
	"BRACE_CLOSE", "BRACE_OPEN", "CONTENT",
	"EMAIL", "ESCAPED_QUOTE",
	"HARD_PAREN_CLOSE", "HARD_PAREN_OPEN",
	"HTML_TAG_OPEN", "HTML_TAG_CLOSE", "HTML_TAG_VOID_CLOSE",
	"IDENTIFIER", "KEYWORD", "LOGICAL",
	"NEWLINE", "NUMERIC_CONTENT", "OPERATOR",
	"PAREN_CLOSE", "PAREN_OPEN", "PERIOD",
	"SINGLE_QUOTE", "DOUBLE_QUOTE", "TEXT_TAG_CLOSE",
	"TEXT_TAG_OPEN", "COMMENT_TAG_OPEN", "COMMENT_TAG_CLOSE",
	"WHITESPACE", "AT_BLOCK", "AT_SECTION"}

type Token struct {
	Text    string
	TypeStr string
	Type    int
	Line    int
	Pos     int
}

func (token Token) debug() {
	textStr := strings.Replace(token.Text, "\n", "\\n", -1)
	textStr = strings.Replace(textStr, "\t", "\\t", -1)
	fmt.Printf("Token: %-20s Location:(%-2d, %d) Value: %s\n",
		token.TypeStr, token.Line, token.Pos, textStr)
}

type TokenMatch struct {
	Type  int
	Text  string
	Regex *regexp.Regexp
}

// The order is important
var TokenMatches = []TokenMatch{
	TokenMatch{EMAIL, "EMAIL", rec(`([a-zA-Z0-9.%]+@[a-zA-Z0-9.\-]+\.(?:ca|co\.uk|com|edu|net|org))\b`)},
	TokenMatch{HTML_TAG_OPEN, "HTML_TAG_OPEN", rec(`(<[a-zA-Z@]+?[^>]*?["a-zA-Z]*>)`)},
	TokenMatch{HTML_TAG_CLOSE, "HTML_TAG_CLOSE", rec(`(</[^>@]+?>)`)},
	TokenMatch{HTML_TAG_VOID_CLOSE, "HTML_TAG_VOID_CLOSE", rec(`(\/\s*>)`)},
	TokenMatch{KEYWORD, "KEYWORD", rec(`(case|do|else|for|func|goto|if|return|switch|var|with)([^\d\w])`)},
	TokenMatch{IDENTIFIER, "IDENTIFIER", rec(`([_$a-zA-Z][_$a-zA-Z0-9]*(\.\.\.)?)`)}, //need verify
	TokenMatch{OPERATOR, "OPERATOR", rec(`(==|!=|>>|<<|>=|<=|>|<|\+|-|\/|\*|\^|%|\:|\?)`)},
	TokenMatch{ESCAPED_QUOTE, "ESCAPED_QUOTE", rec(`(\\+['\"])`)},
	TokenMatch{NUMERIC_CONTENT, "NUMERIC_CONTENT", rec(`([0-9]+)`)},
	TokenMatch{CONTENT, "CONTENT", rec(`([^\s})@.]+?)`)},
}
