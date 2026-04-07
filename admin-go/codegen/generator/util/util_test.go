package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileIfChangedSkipsIdenticalContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "demo.txt")
	if err := os.WriteFile(path, []byte("hello"), 0o644); err != nil {
		t.Fatalf("write seed file: %v", err)
	}

	written, err := WriteFileIfChanged(path, []byte("hello"))
	if err != nil {
		t.Fatalf("WriteFileIfChanged failed: %v", err)
	}
	if written {
		t.Fatal("WriteFileIfChanged should skip identical content")
	}
}

func TestWriteFileIfChangedWritesChangedContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "demo.txt")
	if err := os.WriteFile(path, []byte("old"), 0o644); err != nil {
		t.Fatalf("write seed file: %v", err)
	}

	written, err := WriteFileIfChanged(path, []byte("new"))
	if err != nil {
		t.Fatalf("WriteFileIfChanged failed: %v", err)
	}
	if !written {
		t.Fatal("WriteFileIfChanged should write changed content")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read updated file: %v", err)
	}
	if string(data) != "new" {
		t.Fatalf("unexpected file content: %q", string(data))
	}
}

func TestGenerateFilesSkipsUnchangedWritesInForceMode(t *testing.T) {
	root := t.TempDir()
	tplDir := filepath.Join(root, "tpl")
	outDir := filepath.Join(root, "out")
	if err := os.MkdirAll(tplDir, 0o755); err != nil {
		t.Fatalf("mkdir tpl dir: %v", err)
	}

	tplPath := filepath.Join(tplDir, "demo.tpl")
	if err := os.WriteFile(tplPath, []byte("hello {{.Name}}"), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	outPath := filepath.Join(outDir, "demo", "dept.txt")
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		t.Fatalf("mkdir output dir: %v", err)
	}
	if err := os.WriteFile(outPath, []byte("hello world"), 0o644); err != nil {
		t.Fatalf("write existing output: %v", err)
	}

	generated, err := GenerateFiles([]TemplateMapping{{TplFile: "demo.tpl", OutputPath: "{app}/{module}.txt"}}, tplDir, outDir, "demo", "dept", true, map[string]string{"Name": "world"})
	if err != nil {
		t.Fatalf("GenerateFiles failed: %v", err)
	}
	if len(generated) != 0 {
		t.Fatalf("GenerateFiles should skip unchanged files, got %v", generated)
	}
}

func TestGenerateFilesWritesChangedContentInForceMode(t *testing.T) {
	root := t.TempDir()
	tplDir := filepath.Join(root, "tpl")
	outDir := filepath.Join(root, "out")
	if err := os.MkdirAll(tplDir, 0o755); err != nil {
		t.Fatalf("mkdir tpl dir: %v", err)
	}

	tplPath := filepath.Join(tplDir, "demo.tpl")
	if err := os.WriteFile(tplPath, []byte("hello {{.Name}}"), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	outPath := filepath.Join(outDir, "demo", "dept.txt")
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		t.Fatalf("mkdir output dir: %v", err)
	}
	if err := os.WriteFile(outPath, []byte("hello old"), 0o644); err != nil {
		t.Fatalf("write existing output: %v", err)
	}

	generated, err := GenerateFiles([]TemplateMapping{{TplFile: "demo.tpl", OutputPath: "{app}/{module}.txt"}}, tplDir, outDir, "demo", "dept", true, map[string]string{"Name": "world"})
	if err != nil {
		t.Fatalf("GenerateFiles failed: %v", err)
	}
	if len(generated) != 1 || generated[0] != outPath {
		t.Fatalf("GenerateFiles should report rewritten file, got %v", generated)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read updated output: %v", err)
	}
	if string(data) != "hello world" {
		t.Fatalf("unexpected generated content: %q", string(data))
	}
}
