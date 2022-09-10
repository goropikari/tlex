package generator

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
	buf := bytes.NewBufferString("\n" + ruleStr)
	rules := make([][]string, 0)
	for {
		rule := p.readRule(buf)
		blk := p.readBlock(buf)
		if blk == "" {
			break
		}
		rules = append(rules, []string{rule, blk})
	}

	return rules
}

func (p *Parser) readRule(reader io.RuneReader) string {
	var prev rune
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return ""
			}
			panic(err)
		}
		if prev == '\n' && r == '"' {
			break
		}
		prev = r
	}

	rs := make([]rune, 0)
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}

		switch r {
		case '\\':
			if prev == '\\' {
				rs = append(rs, r)
				prev = 0
				continue
			}
			prev = r
			continue
		case 'n':
			if prev == '\\' {
				rs = append(rs, '\n')
				prev = 0
				continue
			}
		case 'r':
			if prev == '\\' {
				rs = append(rs, '\r')
				prev = 0
				continue
			}
		case 't':
			if prev == '\\' {
				rs = append(rs, '\t')
				prev = 0
				continue
			}
		case '"':
			if prev == '\\' {
				prev = r
				rs = append(rs, r)
				continue
			}

			return string(rs)
		}
		prev = r
		rs = append(rs, r)
	}
	return ""
}

func (p *Parser) readBlock(reader io.RuneReader) string {
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return ""
			}
			panic(err)
		}
		if r == '{' {
			break
		}
	}

	nparen := 1
	rs := []rune{'{'}
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}
		rs = append(rs, r)
		switch r {
		case '{':
			// コメント, 文字列中に { が使われていたときはインクリメントしない処理が本来は必要
			nparen++
		case '}':
			// コメント, 文字列中に { が使われていたときはデクリメントしない処理が本来は必要
			nparen--
			if nparen == 0 {
				return string(rs)
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
