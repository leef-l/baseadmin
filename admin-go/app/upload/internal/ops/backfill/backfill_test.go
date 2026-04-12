package backfill

import (
	"testing"

	"gbaseadmin/app/upload/internal/consts"
)

func TestExtractRelativeDirFromURL(t *testing.T) {
	cases := []struct {
		name string
		url  string
		want string
	}{
		{name: "upload prefix", url: "/upload/2026-04-13/demo.txt", want: "2026-04-13"},
		{name: "resource prefix", url: "/resource/cert/demo.pem", want: "cert"},
		{name: "absolute url", url: "https://cdn.example.com/upload/a/b/demo.png", want: "a/b"},
		{name: "root file", url: "/upload/demo.txt", want: ""},
		{name: "blank", url: " ", want: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := extractRelativeDirFromURL(tc.url); got != tc.want {
				t.Fatalf("extractRelativeDirFromURL(%q) = %q, want %q", tc.url, got, tc.want)
			}
		})
	}
}

func TestCompileRulePathPattern(t *testing.T) {
	regex, score, ok, err := compileRulePathPattern("ci-smoke/{Y-m-d}/{ext}")
	if err != nil {
		t.Fatalf("compileRulePathPattern error: %v", err)
	}
	if !ok {
		t.Fatal("compileRulePathPattern should build matcher")
	}
	if score <= 0 {
		t.Fatalf("compileRulePathPattern invalid score: %d", score)
	}
	if !regex.MatchString("ci-smoke/2026-04-13/txt") {
		t.Fatal("compileRulePathPattern should match rendered path")
	}
	if regex.MatchString("ci-smoke/2026/04/13/txt") {
		t.Fatal("compileRulePathPattern should reject unmatched path")
	}
}

func TestRuleFileTypeMatchesNormalizesDelimiters(t *testing.T) {
	if !ruleFileTypeMatches("png， jpg;gif", "jpg") {
		t.Fatal("ruleFileTypeMatches should normalize mixed delimiters")
	}
}

func TestSelectDirIDForBackfillPrefersTypeRule(t *testing.T) {
	typeRuleRegex, _, _, err := compileRulePathPattern("assets/{Y-m-d}")
	if err != nil {
		t.Fatalf("compile type rule: %v", err)
	}
	defaultRuleRegex, _, _, err := compileRulePathPattern("assets/{Y-m-d}")
	if err != nil {
		t.Fatalf("compile default rule: %v", err)
	}

	matchers := []ruleMatcher{
		{
			dirID:        100,
			category:     consts.DirRuleCategory1,
			storageTypes: "1",
			regex:        defaultRuleRegex,
			score:        10 + categoryBackfillWeight(consts.DirRuleCategory1),
		},
		{
			dirID:        200,
			category:     consts.DirRuleCategory2,
			fileType:     "png,jpg",
			storageTypes: "1",
			regex:        typeRuleRegex,
			score:        10 + categoryBackfillWeight(consts.DirRuleCategory2),
		},
	}

	got, matched, ambiguous := selectDirIDForBackfill(matchers, "assets/2026-04-13", "png", 1)
	if ambiguous {
		t.Fatal("selectDirIDForBackfill should not be ambiguous")
	}
	if !matched {
		t.Fatal("selectDirIDForBackfill should match")
	}
	if got != 200 {
		t.Fatalf("selectDirIDForBackfill = %d, want 200", got)
	}
}

func TestSelectDirIDForBackfillKeepsSourceRuleAmbiguous(t *testing.T) {
	sourceRegex, _, _, err := compileRulePathPattern("shared/{Y-m-d}")
	if err != nil {
		t.Fatalf("compile source rule: %v", err)
	}
	defaultRegex, _, _, err := compileRulePathPattern("shared/{Y-m-d}")
	if err != nil {
		t.Fatalf("compile default rule: %v", err)
	}

	matchers := []ruleMatcher{
		{
			dirID:        100,
			category:     consts.DirRuleCategory1,
			storageTypes: "1",
			regex:        defaultRegex,
			score:        8 + categoryBackfillWeight(consts.DirRuleCategory1),
		},
		{
			dirID:        200,
			category:     consts.DirRuleCategory3,
			storageTypes: "1",
			regex:        sourceRegex,
			score:        8 + categoryBackfillWeight(consts.DirRuleCategory3),
		},
	}

	got, matched, ambiguous := selectDirIDForBackfill(matchers, "shared/2026-04-13", "txt", 1)
	if got != 0 || matched || !ambiguous {
		t.Fatalf("selectDirIDForBackfill ambiguous mismatch: got=%d matched=%t ambiguous=%t", got, matched, ambiguous)
	}
}
