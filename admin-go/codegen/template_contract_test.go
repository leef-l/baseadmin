package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gbaseadmin/codegen/parser"
)

func TestFormTemplateCoversSupportedGeneratedComponents(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(root, "templates", "frontend", "form.tpl"))
	if err != nil {
		t.Fatalf("read form.tpl failed: %v", err)
	}

	formTemplate := string(content)
	handledByFallback := map[string]bool{
		parser.ComponentInput: true,
	}

	for _, componentName := range parser.SupportedComponentNames() {
		if handledByFallback[componentName] {
			continue
		}

		token := `eq .Component "` + componentName + `"`
		if !strings.Contains(formTemplate, token) {
			t.Fatalf("form.tpl missing component branch: %s", componentName)
		}
	}
}
