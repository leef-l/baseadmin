package util

import (
	"os"
	"path/filepath"
	"reflect"
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

func TestPlanAndCommitSkipUnchangedWritesInForceMode(t *testing.T) {
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

	plans, err := PlanFiles([]TemplateMapping{{TplFile: "demo.tpl", OutputPath: "{app}/{module}.txt"}}, tplDir, outDir, "demo", "dept", true, map[string]string{"Name": "world"})
	if err != nil {
		t.Fatalf("PlanFiles failed: %v", err)
	}

	generated, err := CommitPlannedFiles(plans)
	if err != nil {
		t.Fatalf("CommitPlannedFiles failed: %v", err)
	}
	if len(generated) != 0 {
		t.Fatalf("CommitPlannedFiles should skip unchanged files, got %v", generated)
	}
}

func TestPlanAndCommitWritesChangedContentInForceMode(t *testing.T) {
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

	plans, err := PlanFiles([]TemplateMapping{{TplFile: "demo.tpl", OutputPath: "{app}/{module}.txt"}}, tplDir, outDir, "demo", "dept", true, map[string]string{"Name": "world"})
	if err != nil {
		t.Fatalf("PlanFiles failed: %v", err)
	}

	generated, err := CommitPlannedFiles(plans)
	if err != nil {
		t.Fatalf("CommitPlannedFiles failed: %v", err)
	}
	if len(generated) != 1 || generated[0] != outPath {
		t.Fatalf("CommitPlannedFiles should report rewritten file, got %v", generated)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read updated output: %v", err)
	}
	if string(data) != "hello world" {
		t.Fatalf("unexpected generated content: %q", string(data))
	}
}

func TestPlanFilesCapturesCreateSkipAndProtectActions(t *testing.T) {
	root := t.TempDir()
	tplDir := filepath.Join(root, "tpl")
	outDir := filepath.Join(root, "out")
	if err := os.MkdirAll(tplDir, 0o755); err != nil {
		t.Fatalf("mkdir tpl dir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tplDir, "new.tpl"), []byte("new {{.Name}}"), 0o644); err != nil {
		t.Fatalf("write new template: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tplDir, "existing.tpl"), []byte("existing {{.Name}}"), 0o644); err != nil {
		t.Fatalf("write existing template: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tplDir, "enhance.tpl"), []byte("enhance {{.Name}}"), 0o644); err != nil {
		t.Fatalf("write enhance template: %v", err)
	}

	existingPath := filepath.Join(outDir, "demo", "existing.txt")
	if err := os.MkdirAll(filepath.Dir(existingPath), 0o755); err != nil {
		t.Fatalf("mkdir existing dir: %v", err)
	}
	if err := os.WriteFile(existingPath, []byte("existing old"), 0o644); err != nil {
		t.Fatalf("write existing file: %v", err)
	}

	plans, err := PlanFiles([]TemplateMapping{
		{TplFile: "new.tpl", OutputPath: "{app}/new.txt"},
		{TplFile: "existing.tpl", OutputPath: "{app}/existing.txt"},
		{TplFile: "enhance.tpl", OutputPath: "{app}/module_enhance.txt"},
	}, tplDir, outDir, "demo", "dept", false, map[string]string{"Name": "world"})
	if err != nil {
		t.Fatalf("PlanFiles failed: %v", err)
	}

	gotActions := []FilePlanAction{plans[0].Action, plans[1].Action, plans[2].Action}
	wantActions := []FilePlanAction{FilePlanActionCreate, FilePlanActionSkipExisting, FilePlanActionCreate}
	if !reflect.DeepEqual(gotActions, wantActions) {
		t.Fatalf("unexpected plan actions: got=%v want=%v", gotActions, wantActions)
	}

	forcedPlans, err := PlanFiles([]TemplateMapping{
		{TplFile: "enhance.tpl", OutputPath: "{app}/module_enhance.txt"},
	}, tplDir, outDir, "demo", "dept", true, map[string]string{"Name": "world"})
	if err != nil {
		t.Fatalf("forced PlanFiles failed: %v", err)
	}
	if len(forcedPlans) != 1 || forcedPlans[0].Action != FilePlanActionProtectEnhance {
		t.Fatalf("expected protect enhance action, got %+v", forcedPlans)
	}
}

func TestCommitPlannedFilesWritesOnlyCreateAndUpdate(t *testing.T) {
	root := t.TempDir()
	createPath := filepath.Join(root, "create.txt")
	updatePath := filepath.Join(root, "update.txt")
	if err := os.WriteFile(updatePath, []byte("old"), 0o644); err != nil {
		t.Fatalf("write seed update file: %v", err)
	}

	plans := []PlannedFile{
		{OutputPath: createPath, Action: FilePlanActionCreate, Content: []byte("create"), Bytes: 6},
		{OutputPath: updatePath, Action: FilePlanActionUpdate, Content: []byte("update"), Bytes: 6},
		{OutputPath: filepath.Join(root, "skip.txt"), Action: FilePlanActionSkipExisting, Content: []byte("skip"), Bytes: 4},
	}

	generated, err := CommitPlannedFiles(plans)
	if err != nil {
		t.Fatalf("CommitPlannedFiles failed: %v", err)
	}
	if len(generated) != 2 {
		t.Fatalf("unexpected generated files: %v", generated)
	}

	data, err := os.ReadFile(createPath)
	if err != nil || string(data) != "create" {
		t.Fatalf("create file mismatch: %v %q", err, string(data))
	}
	data, err = os.ReadFile(updatePath)
	if err != nil || string(data) != "update" {
		t.Fatalf("update file mismatch: %v %q", err, string(data))
	}
	if _, err := os.Stat(filepath.Join(root, "skip.txt")); !os.IsNotExist(err) {
		t.Fatalf("skip file should not be written, stat err=%v", err)
	}
}
