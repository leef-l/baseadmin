package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"gbaseadmin/codegen/generator/util"
	"gbaseadmin/codegen/internal/verifytemplates"
	"gbaseadmin/codegen/parser"
)

func loadTemplateCaseMeta(t *testing.T, key string) *parser.TableMeta {
	t.Helper()

	for _, tc := range verifytemplates.Cases() {
		if tc.Key == key {
			meta := *tc.Meta
			return &meta
		}
	}

	t.Fatalf("template case %q not found", key)
	return nil
}

func renderFrontendTemplate(t *testing.T, tplFile string, meta *parser.TableMeta) string {
	t.Helper()

	root, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}

	tplPath := filepath.Join(root, "templates", tplFile)
	tpl, err := template.New(filepath.Base(tplPath)).Funcs(util.SharedFuncMap).ParseFiles(tplPath)
	if err != nil {
		t.Fatalf("parse template %s failed: %v", tplFile, err)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, meta); err != nil {
		t.Fatalf("render template %s failed: %v", tplFile, err)
	}

	return buf.String()
}

func TestFrontendTemplatesImportDictAPIByDefault(t *testing.T) {
	meta := loadTemplateCaseMeta(t, "demo_article")

	for _, tplFile := range []string{"frontend/form.tpl", "frontend/list.tpl"} {
		output := renderFrontendTemplate(t, tplFile, meta)
		if !strings.Contains(output, "import { getDictByType } from '#/api/system/dict';") {
			t.Fatalf("%s should import the dict api when allow_missing_dict_module is disabled", tplFile)
		}
		if strings.Contains(output, "async function getDictByType") {
			t.Fatalf("%s should not inline the dict fallback when allow_missing_dict_module is disabled", tplFile)
		}
	}
}

func TestFrontendTemplatesInlineDictFallbackWhenAllowed(t *testing.T) {
	meta := loadTemplateCaseMeta(t, "demo_article")
	meta.AllowMissingDictModule = true

	for _, tplFile := range []string{"frontend/form.tpl", "frontend/list.tpl"} {
		output := renderFrontendTemplate(t, tplFile, meta)
		if strings.Contains(output, "#/api/system/dict") {
			t.Fatalf("%s should not import the dict api when allow_missing_dict_module is enabled", tplFile)
		}
		if !strings.Contains(output, "async function getDictByType") {
			t.Fatalf("%s should inline a dict fallback when allow_missing_dict_module is enabled", tplFile)
		}
	}
}
