package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"gbaseadmin/codegen/internal/verifytemplates"
)

var updateGolden = flag.Bool("update-golden", false, "update codegen golden snapshots")

func TestTemplateGoldenSnapshots(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	tplDir := filepath.Join(root, "templates")
	goldenRoot := filepath.Join(root, "testdata", "golden")

	for _, tc := range verifytemplates.Cases() {
		tc := tc
		t.Run(tc.Key, func(t *testing.T) {
			for _, tplFile := range verifytemplates.AllTemplateFiles() {
				tplFile := tplFile
				t.Run(tplFile, func(t *testing.T) {
					rendered, err := verifytemplates.RenderTemplateCase(tplDir, tplFile, tc)
					if err != nil {
						t.Fatalf("renderTemplateCase failed: %v", err)
					}

					goldenPath := filepath.Join(goldenRoot, tc.Key, tplFile+".golden")
					if *updateGolden {
						if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
							t.Fatalf("mkdir golden dir failed: %v", err)
						}
						if err := os.WriteFile(goldenPath, []byte(rendered.Output), 0o644); err != nil {
							t.Fatalf("write golden failed: %v", err)
						}
						return
					}

					want, err := os.ReadFile(goldenPath)
					if err != nil {
						t.Fatalf("read golden failed: %v\nrun: go test ./... -run TestTemplateGoldenSnapshots -args -update-golden", err)
					}
					if rendered.Output != string(want) {
						t.Fatalf("golden mismatch: %s\nrun: go test ./... -run TestTemplateGoldenSnapshots -args -update-golden", goldenPath)
					}
				})
			}
		})
	}
}
