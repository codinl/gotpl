package gotpl

import (
	"fmt"
	"regexp"
	"strings"
)

type Lexer struct {
	Text    string
	Matches []TokenMatch
}

func (lexer *Lexer) Scan() ([]Token, error) {
	tokens := []Token{}
	text := strings.Replace(lexer.Text, "\r\n", "\n", -1)
	text = strings.Replace(lexer.Text, "\r", "\n", -1)
	text += "\n"
	cur, line, pos := 0, 0, 0
	for cur < len(text) {
		val, left := text[cur], text[cur:]
		var token Token
		switch val {
		case '\n':
			token = makeToken(string(val), NEWLINE)
		case ' ', '\t', '\f', '\r':
			token = makeToken(string(val), WHITESPACE)
		case '(':
			token = makeToken(string(val), PAREN_OPEN)
		case ')':
			token = makeToken(string(val), PAREN_CLOSE)
		case '[':
			token = makeToken(string(val), HARD_PAREN_OPEN)
		case ']':
			token = makeToken(string(val), HARD_PAREN_CLOSE)
		case '{':
			token = makeToken(string(val), BRACE_OPEN)
		case '}':
			token = makeToken(string(val), BRACE_CLOSE)
		case '"', '`':
			token = makeToken(string(val), DOUBLE_QUOTE)
		case '\'':
			token = makeToken(string(val), SINGLE_QUOTE)
		case '.':
			token = makeToken(string(val), PERIOD)
		case '@':
			if peekNext(":", left[1:]) {
				token = makeToken("@:", AT_COLON)
			} else if peekNext("*", left[1:]) {
				token = makeToken("@*", AT_STAR_OPEN)
			} else if peekNext("block", left[1:]) {
				token = makeToken("@block", AT_BLOCK)
			} else if peekNext("section", left[1:]) {
				token = makeToken("@section", AT_SECTION)
			} else {
				token = makeToken("@", AT)
			}
		default:
			if peekNext("*@", left) {
				token = makeToken("*@", AT_STAR_CLOSE)
			} else if peekNext("<text>", left) {
				token = makeToken("<text>", TEXT_TAG_OPEN)
			} else if peekNext("</text>", left) {
				token = makeToken("</text>", TEXT_TAG_CLOSE)
			} else if peekNext("<!--", left) {
				token = makeToken("<!--", COMMENT_TAG_OPEN)
			} else if peekNext("-->", left) {
				token = makeToken("-->", COMMENT_TAG_CLOSE)
			} else {
				//try rec
				match := false
				for _, m := range lexer.Matches {
					found := m.Regex.FindIndex([]byte(left))
					if found != nil {
						match = true
						tokenVal := left[found[0]:found[1]]
						if m.Type == HTML_TAG_OPEN {
							tokenVal = tagClean(tokenVal)
						} else if m.Type == KEYWORD {
							tokenVal = keyClean(tokenVal)
						}
						token = makeToken(tokenVal, m.Type)
						break
					}
				}
				if !match {
					return tokens, fmt.Errorf("%d:%d: Illegal character: %s",
						line, pos, string(text[pos]))
				}
			}
		}
		token.Line, token.Pos = line, pos
		tokens = append(tokens, token)
		cur += len(token.Text)
		if token.Type == NEWLINE {
			line, pos = line+1, 0
		} else {
			pos += len(token.Text)
		}
	}
	return tokens, nil
}

// Why we need this: Go's regexp DO NOT support lookahead assertion
func regRemoveTail(text string, regs []string) string {
	res := text
	for _, reg := range regs {
		regc := regexp.MustCompile(reg)
		found := regc.FindIndex([]byte(res))
		if found != nil {
			res = res[:found[0]] //BUG?
		}
	}
	return res
}

func tagClean(text string) string {
	regs := []string{
		`([a-zA-Z0-9.%]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4})\b`,
		`(@)`,
		`(\/\s*>)`}
	return regRemoveTail(text, regs)
}

func keyClean(text string) string {
	pos := len(text) - 1
	for {
		v := text[pos]
		if (v >= 'a' && v <= 'z') ||
			(v >= 'A' && v <= 'Z') {
			break
		} else {
			pos--
		}
	}
	return text[:pos+1]
}

func peekNext(expect string, text string) bool {
	if strings.HasPrefix(text, expect) {
		return true
	}
	return false
}

func makeToken(val string, tokenType int) Token {
	return Token{val, TokenTypeNames[tokenType], tokenType, 0, 0}
}
