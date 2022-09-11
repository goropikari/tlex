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
	buf := bytes.NewBufferString(ruleStr)
	rules := make([][]string, 0)
	for {
		if err := skipWhitespace(buf); err != nil {
			if errors.Is(err, io.EOF) {
				return rules
			}
			panic(err)
		}
		rule := p.readRule(buf)
		blk := p.readBlock(buf)
		if blk == "" {
			break
		}
		rules = append(rules, []string{rule, blk})
	}

	return rules
}

func skipWhitespace(reader io.RuneScanner) error {
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return err
		}
		switch r {
		case '\n', '\r', '\t', ' ':
			continue
		}
		if err := reader.UnreadRune(); err != nil {
			return err
		}
		break
	}

	return nil
}

func (p *Parser) readRule(reader io.RuneScanner) string {
	inRange := false
	rs := make([]rune, 0)
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			panic(err)
		}

		switch r {
		case '"':
			rs = append(rs, r)
			for {
				r, _, err = reader.ReadRune()
				if err != nil {
					panic(err)
				}

				switch r {
				case '\\':
					nr, err := nextRune(reader)
					if err != nil {
						panic(err)
					}
					reader.ReadRune()
					if nr == '"' {
						rs = append(rs, '\\')
						rs = append(rs, '"')
					}
				case '"':
					rs = append(rs, r)
					return string(rs)
				default:
					rs = append(rs, r)
				}
			}
		case '[':
			inRange = true
		case ']':
			inRange = false
		case ' ':
			if !inRange {
				return string(rs)
			}
		case '\t':
			return string(rs)
		}

		rs = append(rs, r)
	}
}

func nextRune(reader io.RuneScanner) (rune, error) {
	r, _, err := reader.ReadRune()
	if err != nil {
		return 0, err
	}
	if err := reader.UnreadRune(); err != nil {
		return 0, err
	}
	return r, nil
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

	_ = p.readUntil(p.r, "%%\n")
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
