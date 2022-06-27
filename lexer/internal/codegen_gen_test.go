
// Code generated by Participle. DO NOT EDIT.
package internal_test

import (
	"io"
	"strings"
	"unicode/utf8"
	"regexp/syntax"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var _ syntax.Op

var Lexer lexer.Definition = definitionImpl{}

type definitionImpl struct {}

func (definitionImpl) Symbols() map[string]lexer.TokenType {
	return map[string]lexer.TokenType{
      "Char": -11,
      "EOF": -1,
      "Escaped": -8,
      "Expr": -10,
      "ExprEnd": -6,
      "Ident": -5,
      "Oper": -4,
      "String": -7,
      "StringEnd": -9,
      "Whitespace": -3,
	}
}

func (definitionImpl) LexString(filename string, s string) (lexer.Lexer, error) {
	return &lexerImpl{
		s: s,
		pos: lexer.Position{
			Filename: filename,
			Line:     1,
			Column:   1,
		},
		states: []lexerState{lexerState{name: "Root"}},
	}, nil
}

func (d definitionImpl) LexBytes(filename string, b []byte) (lexer.Lexer, error) {
	return d.LexString(filename, string(b))
}

func (d definitionImpl) Lex(filename string, r io.Reader) (lexer.Lexer, error) {
	s := &strings.Builder{}
	_, err := io.Copy(s, r)
	if err != nil {
		return nil, err
	}
	return d.LexString(filename, s.String())
}

type lexerState struct {
	name    string
	groups  []string
}

type lexerImpl struct {
	s       string
	p       int
	pos     lexer.Position
	states  []lexerState
}

func (l *lexerImpl) Next() (lexer.Token, error) {
	if l.p == len(l.s) {
		return lexer.EOFToken(l.pos), nil
	}
	var (
		state = l.states[len(l.states)-1]
		groups []int
		sym lexer.TokenType
	)
	switch state.name {
	case "Expr":if match := matchString(l.s, l.p); match[1] != 0 {
			sym = -7
			groups = match[:]
			l.states = append(l.states, lexerState{name: "String"})
		} else if match := matchWhitespace(l.s, l.p); match[1] != 0 {
			sym = -3
			groups = match[:]
		} else if match := matchOper(l.s, l.p); match[1] != 0 {
			sym = -4
			groups = match[:]
		} else if match := matchIdent(l.s, l.p); match[1] != 0 {
			sym = -5
			groups = match[:]
		} else if match := matchExprEnd(l.s, l.p); match[1] != 0 {
			sym = -6
			groups = match[:]
			l.states = l.states[:len(l.states)-1]
		}
	case "Root":if match := matchString(l.s, l.p); match[1] != 0 {
			sym = -7
			groups = match[:]
			l.states = append(l.states, lexerState{name: "String"})
		}
	case "String":if match := matchEscaped(l.s, l.p); match[1] != 0 {
			sym = -8
			groups = match[:]
		} else if match := matchStringEnd(l.s, l.p); match[1] != 0 {
			sym = -9
			groups = match[:]
			l.states = l.states[:len(l.states)-1]
		} else if match := matchExpr(l.s, l.p); match[1] != 0 {
			sym = -10
			groups = match[:]
			l.states = append(l.states, lexerState{name: "Expr"})
		} else if match := matchChar(l.s, l.p); match[1] != 0 {
			sym = -11
			groups = match[:]
		}
	}
	if groups == nil {
		sample := []rune(l.s[l.p:])
		if len(sample) > 16 {
			sample = append(sample[:16], []rune("...")...)
		}
		return lexer.Token{}, participle.Errorf(l.pos, "invalid input text %q", sample)
	}
	pos := l.pos
	span := l.s[groups[0]:groups[1]]
	l.p = groups[1]
	l.pos.Advance(span)
	return lexer.Token{
		Type:  sym,
		Value: span,
		Pos:   pos,
	}, nil
}

func (l *lexerImpl) sgroups(match []int) []string {
	sgroups := make([]string, len(match)/2)
	for i := 0; i < len(match)-1; i += 2 {
		sgroups[i/2] = l.s[l.p+match[i]:l.p+match[i+1]]
	}
	return sgroups
}


// "
func matchString(s string, p int) (groups [2]int) {
if p < len(s) && s[p] == '"' {
groups[0] = p
groups[1] = p + 1
}
return
}

// [\t-\n\f-\r ]+
func matchWhitespace(s string, p int) (groups [2]int) {
// [\t-\n\f-\r ] (CharClass)
l0 := func(s string, p int) int {
if len(s) <= p { return -1 }
rn := s[p]
switch {
case rn >= '\t' && rn <= '\n': return p+1
case rn >= '\f' && rn <= '\r': return p+1
case rn == ' ': return p+1
}
return -1
}
// [\t-\n\f-\r ]+ (Plus)
l1 := func(s string, p int) int {
if p = l0(s, p); p == -1 { return -1 }
for len(s) > p {
if np := l0(s, p); np == -1 { return p } else { p = np }
}
return p
}
np := l1(s, p)
if np == -1 {
  return
}
groups[0] = p
groups[1] = np
return
}

// [%\*-\+\-/]
func matchOper(s string, p int) (groups [2]int) {
// [%\*-\+\-/] (CharClass)
l0 := func(s string, p int) int {
if len(s) <= p { return -1 }
rn := s[p]
switch {
case rn == '%': return p+1
case rn >= '*' && rn <= '+': return p+1
case rn == '-': return p+1
case rn == '/': return p+1
}
return -1
}
np := l0(s, p)
if np == -1 {
  return
}
groups[0] = p
groups[1] = np
return
}

// [0-9A-Z_a-z]+
func matchIdent(s string, p int) (groups [2]int) {
// [0-9A-Z_a-z] (CharClass)
l0 := func(s string, p int) int {
if len(s) <= p { return -1 }
rn := s[p]
switch {
case rn >= '0' && rn <= '9': return p+1
case rn >= 'A' && rn <= 'Z': return p+1
case rn == '_': return p+1
case rn >= 'a' && rn <= 'z': return p+1
}
return -1
}
// [0-9A-Z_a-z]+ (Plus)
l1 := func(s string, p int) int {
if p = l0(s, p); p == -1 { return -1 }
for len(s) > p {
if np := l0(s, p); np == -1 { return p } else { p = np }
}
return p
}
np := l1(s, p)
if np == -1 {
  return
}
groups[0] = p
groups[1] = np
return
}

// \}
func matchExprEnd(s string, p int) (groups [2]int) {
if p < len(s) && s[p] == '}' {
groups[0] = p
groups[1] = p + 1
}
return
}

// \\(?-s:.)
func matchEscaped(s string, p int) (groups [2]int) {
// \\ (Literal)
l0 := func(s string, p int) int {
if p < len(s) && s[p] == '\\' { return p+1 }
return -1
}
// (?-s:.) (AnyCharNotNL)
l1 := func(s string, p int) int {
var (rn rune; n int)
if s[p] < utf8.RuneSelf {
  rn, n = rune(s[p]), 1
} else {
  rn, n = utf8.DecodeRuneInString(s[p:])
}
if len(s) <= p+n || rn == '\n' { return -1 }
return p+n
}
// \\(?-s:.) (Concat)
l2 := func(s string, p int) int {
if p = l0(s, p); p == -1 { return -1 }
if p = l1(s, p); p == -1 { return -1 }
return p
}
np := l2(s, p)
if np == -1 {
  return
}
groups[0] = p
groups[1] = np
return
}

// "
func matchStringEnd(s string, p int) (groups [2]int) {
if p < len(s) && s[p] == '"' {
groups[0] = p
groups[1] = p + 1
}
return
}

// \$\{
func matchExpr(s string, p int) (groups [2]int) {
if p+2 < len(s) && s[p:p+2] == "${" {
groups[0] = p
groups[1] = p + 2
}
return
}

// [^"\$\\]+
func matchChar(s string, p int) (groups [2]int) {
// [^"\$\\] (CharClass)
l0 := func(s string, p int) int {
if len(s) <= p { return -1 }
var (rn rune; n int)
if s[p] < utf8.RuneSelf {
  rn, n = rune(s[p]), 1
} else {
  rn, n = utf8.DecodeRuneInString(s[p:])
}
switch {
case rn >= '\x00' && rn <= '!': return p+1
case rn == '#': return p+1
case rn >= '%' && rn <= '[': return p+1
case rn >= ']' && rn <= '\U0010ffff': return p+n
}
return -1
}
// [^"\$\\]+ (Plus)
l1 := func(s string, p int) int {
if p = l0(s, p); p == -1 { return -1 }
for len(s) > p {
if np := l0(s, p); np == -1 { return p } else { p = np }
}
return p
}
np := l1(s, p)
if np == -1 {
  return
}
groups[0] = p
groups[1] = np
return
}
