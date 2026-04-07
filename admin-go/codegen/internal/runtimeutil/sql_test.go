package runtimeutil

import (
	"reflect"
	"testing"
)

func TestSplitSQLStatementsBasic(t *testing.T) {
	sqlText := `
-- comment
CREATE TABLE demo_a (
  id BIGINT
);

-- another comment
INSERT INTO demo_a VALUES (1);

`

	got := SplitSQLStatements(sqlText)
	want := []string{
		"CREATE TABLE demo_a (\n  id BIGINT\n)",
		"INSERT INTO demo_a VALUES (1)",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SplitSQLStatements mismatch:\nwant=%q\ngot=%q", want, got)
	}
}

func TestSplitSQLStatementsHandlesQuotedSemicolonsAndComments(t *testing.T) {
	sqlText := `
CREATE TABLE demo_b (
  note VARCHAR(64) DEFAULT 'hello;world'
);
INSERT INTO demo_b VALUES ('a;1', "b;2"); -- trailing comment
/* block comment with ; */
INSERT INTO demo_b VALUES ('/* literal */', 'line -- keep ; here');
`

	got := SplitSQLStatements(sqlText)
	want := []string{
		"CREATE TABLE demo_b (\n  note VARCHAR(64) DEFAULT 'hello;world'\n)",
		`INSERT INTO demo_b VALUES ('a;1', "b;2")`,
		"INSERT INTO demo_b VALUES ('/* literal */', 'line -- keep ; here')",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SplitSQLStatements complex mismatch:\nwant=%q\ngot=%q", want, got)
	}
}

func TestSplitSQLStatementsSupportsHashCommentAndEOFLineComment(t *testing.T) {
	sqlText := "CREATE TABLE demo_c (id BIGINT); # keep this ignored\nINSERT INTO demo_c VALUES (1) -- eof comment"

	got := SplitSQLStatements(sqlText)
	want := []string{
		"CREATE TABLE demo_c (id BIGINT)",
		"INSERT INTO demo_c VALUES (1)",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SplitSQLStatements hash/eof mismatch:\nwant=%q\ngot=%q", want, got)
	}
}
