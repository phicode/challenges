package parser

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"text/scanner"
)

type Parser struct {
	err     error
	scanner scanner.Scanner
}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) SetString(s string) {
	p.scanner.Init(bytes.NewBufferString(s))
}
func (p *Parser) Error() error {
	return p.err
}

// AcceptString a specific string token
// The Error is set if the token did not match
func (p *Parser) AcceptString(str string) bool {
	if p.err != nil {
		return false
	}
	if p.scanner.Scan() == scanner.EOF {
		p.err = io.EOF
		return false
	}
	text := p.scanner.TokenText()
	if text == str {
		return true
	}
	p.err = fmt.Errorf("%v - expected token %q, got %q", p.scanner.Pos(), str, text)
	return false
}
func (p *Parser) Accept(r rune) bool {
	if p.err != nil {
		return false
	}
	token := p.scanner.Scan()
	if token == scanner.EOF {
		p.err = io.EOF
		return false
	}
	if token == r {
		return true
	}
	p.err = fmt.Errorf("%v - expected token %c(%d), got %c(%d)", p.scanner.Pos(), r, r, token, token)
	return false
}

func (p *Parser) Int() (int, bool) {
	if p.err != nil {
		return 0, false
	}
	if p.scanner.Scan() == scanner.EOF {
		p.err = io.EOF
		return 0, false
	}
	text := p.scanner.TokenText()
	value, err := strconv.Atoi(text)
	if err == nil {
		return value, true
	}
	p.err = fmt.Errorf("%v - %w", p.scanner.Pos(), err)
	return 0, false
}

func (p *Parser) String() (string, bool) {
	if p.err != nil {
		return "", false
	}
	if p.scanner.Scan() == scanner.EOF {
		p.err = io.EOF
		return "", false
	}
	text := p.scanner.TokenText()
	return text, true
}

func (p *Parser) Peek() (rune, bool) {
	if p.err != nil {
		return 0, false
	}
	token := p.scanner.Peek()
	return token, token != scanner.EOF
}
