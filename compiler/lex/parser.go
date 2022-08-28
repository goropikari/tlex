package lex

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

// %{
// // Definitions
// %}
//
// %%
// Rules
// %%
//
// User code section

type Parser struct {
	r *bufio.Reader
}

func NewParser(r *bufio.Reader) *Parser {
	return &Parser{
		r: r,
	}
}

func (p *Parser) Parse() (def string, rules [][]string, userCode string) {
	def, rulesStr, userCode := p.Split()

	return def, p.parseRules(rulesStr), userCode
}

func (p *Parser) parseRules(ruleStr string) [][]string {
	r := bytes.NewBufferString("\n" + ruleStr)
	rules := make([][]string, 0)
	for {
		rule := p.readRule(r)
		blk := p.readBlock(r)
		if blk == "" {
			break
		}
		rules = append(rules, []string{rule, blk})
	}

	return rules
}

func (p *Parser) readRule(r io.RuneReader) string {
	var prev rune
	for {
		ru, _, err := r.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return ""
			}
			panic(err)
		}
		if prev == '\n' && ru == '"' {
			break
		}
		prev = ru
	}

	runes := make([]rune, 0)
	for {
		ru, _, err := r.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}

		switch ru {
		case '\\':
			if prev == '\\' {
				runes = append(runes, ru)
				prev = 0
				continue
			}
			prev = ru
			continue
		case 'n':
			if prev == '\\' {
				runes = append(runes, '\n')
				prev = 0
				continue
			}
		case 'r':
			if prev == '\\' {
				runes = append(runes, '\r')
				prev = 0
				continue
			}
		case 't':
			if prev == '\\' {
				runes = append(runes, '\t')
				prev = 0
				continue
			}
		case '"':
			if prev == '\\' {
				prev = ru
				runes = append(runes, ru)
				continue
			}

			return string(runes)
		}
		prev = ru
		runes = append(runes, ru)
	}
	return ""
}

func (p *Parser) readBlock(r io.RuneReader) string {
	for {
		ru, _, err := r.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return ""
			}
			panic(err)
		}
		if ru == '{' {
			break
		}
	}

	nparen := 1
	runes := []rune{'{'}
	for {
		ru, _, err := r.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}
		runes = append(runes, ru)
		switch ru {
		case '{':
			// コメント, 文字列中に { が使われていたときはインクリメントしない処理が本来は必要
			nparen++
		case '}':
			// コメント, 文字列中に { が使われていたときはデクリメントしない処理が本来は必要
			nparen--
			if nparen == 0 {
				return string(runes)
			}
		}
	}

	return ""
}

func (p *Parser) Split() (def string, rules string, userCode string) {

	_ = p.readUntil(p.r, "%{\n")

	def = p.readUntil(p.r, "%}\n")

	p.readUntil(p.r, "%%\n")
	rules = p.readUntil(p.r, "%%\n")

	var buf bytes.Buffer
	io.Copy(&buf, p.r)
	userCode = buf.String()

	return
}

func (p *Parser) readUntil(r *bufio.Reader, delim string) string {
	var buf bytes.Buffer
	for {
		s, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			panic(errors.New("error"))
		}
		if s == delim {
			return buf.String()
		}
		buf.WriteString(s)
	}

	return ""
}
