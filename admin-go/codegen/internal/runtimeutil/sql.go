package runtimeutil

import "strings"

// SplitSQLStatements 按 SQL 语句边界切分文本，忽略注释并保留字符串中的分号。
func SplitSQLStatements(sqlText string) []string {
	var (
		stmts          []string
		buf            strings.Builder
		inSingleQuote  bool
		inDoubleQuote  bool
		inBacktick     bool
		inLineComment  bool
		inBlockComment bool
	)

	flush := func() {
		stmt := strings.TrimSpace(buf.String())
		if stmt != "" {
			stmts = append(stmts, stmt)
		}
		buf.Reset()
	}

	for i := 0; i < len(sqlText); i++ {
		ch := sqlText[i]
		next := byte(0)
		if i+1 < len(sqlText) {
			next = sqlText[i+1]
		}

		if inLineComment {
			if ch == '\n' {
				inLineComment = false
				buf.WriteByte(ch)
			}
			continue
		}
		if inBlockComment {
			if ch == '*' && next == '/' {
				inBlockComment = false
				i++
			}
			continue
		}

		switch {
		case inSingleQuote:
			buf.WriteByte(ch)
			if ch == '\\' && i+1 < len(sqlText) {
				i++
				buf.WriteByte(sqlText[i])
				continue
			}
			if ch == '\'' {
				if next == '\'' {
					i++
					buf.WriteByte(next)
					continue
				}
				inSingleQuote = false
			}
			continue
		case inDoubleQuote:
			buf.WriteByte(ch)
			if ch == '\\' && i+1 < len(sqlText) {
				i++
				buf.WriteByte(sqlText[i])
				continue
			}
			if ch == '"' {
				if next == '"' {
					i++
					buf.WriteByte(next)
					continue
				}
				inDoubleQuote = false
			}
			continue
		case inBacktick:
			buf.WriteByte(ch)
			if ch == '`' {
				inBacktick = false
			}
			continue
		}

		if startsLineComment(sqlText, i) {
			inLineComment = true
			i++
			continue
		}
		if ch == '#' {
			inLineComment = true
			continue
		}
		if ch == '/' && next == '*' {
			inBlockComment = true
			i++
			continue
		}

		switch ch {
		case '\'':
			inSingleQuote = true
			buf.WriteByte(ch)
		case '"':
			inDoubleQuote = true
			buf.WriteByte(ch)
		case '`':
			inBacktick = true
			buf.WriteByte(ch)
		case ';':
			flush()
		default:
			buf.WriteByte(ch)
		}
	}

	flush()
	return stmts
}

func startsLineComment(sqlText string, idx int) bool {
	if idx+1 >= len(sqlText) || sqlText[idx] != '-' || sqlText[idx+1] != '-' {
		return false
	}
	if idx == 0 {
		return true
	}
	switch sqlText[idx-1] {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}
