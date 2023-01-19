package tokenize

import (
	"fmt"
	"strconv"
	"unicode"
)

var userInput []rune

var currentPos *Position

var singleOpSymbols []string
var compositeOpSymbols []string

func init() {
	//currentPos = NewPosition(1, 0, 0)
	singleOpSymbols = []string{
		"(", ")", "[", "]", "{", "}",
		".", ",", ":", ";",
		"+", "-", "*", "/", "%",
		">", "<",
		"=", "!",
	}
	compositeOpSymbols = []string{
		"==", "!=", ">=", "<=",
		"+=", "-=", "*=", "/=", "%=",
		"&&", "||", ":=",
	}
}

func startWith(s string) bool {
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if len(userInput) <= currentPos.Wat+i || userInput[currentPos.Wat+i] != runes[i] {
			return false
		}
	}
	return true
}

func isIdentRune(r rune) bool {
	return ('a' <= r && r <= 'z') ||
		('A' <= r && r <= 'Z') ||
		('0' <= r && r <= '9') ||
		('_' == r || '.' == r)
}

func isNotEof() bool {
	return currentPos.Wat < len(userInput)
}

func consumeComment() string {
	// "//"
	currentPos.Lat += 2
	currentPos.Wat += 2
	var s string
	for isNotEof() {
		if userInput[currentPos.Wat] == '\n' {
			break
		}
		s += string(userInput[currentPos.Wat])
		currentPos.Lat++
		currentPos.Wat++
	}
	return s
}

func consumeIdent() string {
	var s string
	for isNotEof() {
		if !isIdentRune(userInput[currentPos.Wat]) {
			break
		}
	}
	// textは後で検査
	return s
}

func consumeString() string {
	var s string
	// "
	currentPos.Lat++
	currentPos.Wat++

	for isNotEof() {
		if userInput[currentPos.Wat] == '"' {
			break
		}
		// escaped double quotation
		if userInput[currentPos.Wat] == '\\' && userInput[currentPos.Wat+1] == '"' {
			s += "\""
			currentPos.Lat += 2
			currentPos.Wat += 2
			continue
		}
		// newline
		if userInput[currentPos.Wat] == '\\' && userInput[currentPos.Wat+1] == 'n' {
			s += "\n"
			currentPos.Lat += 2
			currentPos.Wat += 2
			continue
		}
		// tab
		if userInput[currentPos.Wat] == '\\' && userInput[currentPos.Wat+1] == 't' {
			s += "\t"
			currentPos.Lat += 2
			currentPos.Wat += 2
			continue
		}
		// escaped single quotation
		if userInput[currentPos.Wat] == '\\' && userInput[currentPos.Wat+1] == '\'' {
			s += "'"
			currentPos.Lat += 2
			currentPos.Wat += 2
			continue
		}
		// escaped? slash
		if userInput[currentPos.Wat] == '\\' && userInput[currentPos.Wat+1] == '\\' {
			s += "\\"
			currentPos.Lat += 2
			currentPos.Wat += 2
			continue
		}

		s += string(userInput[currentPos.Wat])
		currentPos.Lat++
		currentPos.Wat++
	}

	// "
	currentPos.Lat++
	currentPos.Wat++
	return s
}

func consumeNumber() (string, bool) {
	isFloat := false
	var s string
	for isNotEof() {
		if unicode.IsDigit(userInput[currentPos.Wat]) {
			s += string(userInput[currentPos.Wat])
			currentPos.Lat++
			currentPos.Wat++
			continue
		} else if userInput[currentPos.Wat] == '.' {
			// ポイントの次が、数字じゃなければ強制終了
			if len(userInput) <= currentPos.Wat+1 ||
				!unicode.IsDigit(userInput[currentPos.Wat+1]) {
				break
			}
			s += string(userInput[currentPos.Wat])
			currentPos.Lat++
			currentPos.Wat++
			isFloat = true
			continue
		} else {
			break
		}
	}
	return s, isFloat
}

func consumeWhite() string {
	var s string
	for isNotEof() {
		if userInput[currentPos.Wat] == ' ' || userInput[currentPos.Wat] == '\t' {
			s += string(userInput[currentPos.Wat])
			currentPos.Lat++
			currentPos.Wat++
		} else {
			break
		}
	}
	return s
}

func Tokenize(input string) (*Token, error) {
	userInput = []rune(input)
	currentPos = NewPosition(1, 0, 0)
	var head Token
	cur := &head
inputLoop:
	for isNotEof() {
		// white
		if userInput[currentPos.Wat] == ' ' || userInput[currentPos.Wat] == '\t' {
			// もし、入力の完全な復元をするのであれば、これは捨てるべきではない
			_ = currentPos.Clone()
			_ = consumeWhite()
			continue
		}
		// newline
		if userInput[currentPos.Wat] == '\n' {
			//cur = NewNewlineChain(cur, currentPos.Clone(), "\n")
			currentPos.LineNo++
			currentPos.Lat = 0
			currentPos.Wat++
			continue
		}
		// comment
		if userInput[currentPos.Wat] == '\\' && userInput[currentPos.Wat+1] == '\\' {
			pos := currentPos.Clone()
			s := consumeComment()
			cur = NewCommentChain(cur, pos, s)
			continue
		}
		// op, symbols
		for _, r := range append(compositeOpSymbols, singleOpSymbols...) {
			if startWith(r) {
				cur = NewOpSymbolChain(cur, currentPos.Clone(), r)
				currentPos.Lat += len(r)
				currentPos.Wat += len(r)
				// 直前のループではなく全体をコンティニュー
				continue inputLoop
			}
		}
		// ident
		if isIdentRune(userInput[currentPos.Wat]) && !unicode.IsDigit(userInput[currentPos.Wat]) {
			pos := currentPos.Clone()
			id := consumeIdent()
			cur = NewIdentChain(cur, pos, id)
			continue
		}
		// string
		if userInput[currentPos.Wat] == '"' {
			pos := currentPos.Clone()
			s := consumeString()
			cur = NewLiteralChain(cur, pos, NewStringLiteral(s))
			continue
		}
		// number
		if unicode.IsDigit(userInput[currentPos.Wat]) {
			pos := currentPos.Clone()
			numS, isFloat := consumeNumber()
			if isFloat {
				n, err := strconv.ParseFloat(numS, 64)
				if err != nil {
					return nil, err
				}
				cur = NewLiteralChain(cur, pos, NewFloatLiteral(n))
				continue
			} else {
				n, err := strconv.ParseInt(numS, 10, 0)
				if err != nil {
					return nil, err
				}
				cur = NewLiteralChain(cur, pos, NewIntLiteral(int(n)))
				continue
			}
		}
		return nil, fmt.Errorf("unexpected charactor: %v", userInput[currentPos.Wat])
	}
	cur = NewEofChain(cur, currentPos.Clone())
	return head.Next, nil
}
