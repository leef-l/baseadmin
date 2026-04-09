package parser

import "testing"

func TestSupportedComponentNamesAreUniqueAndRegistered(t *testing.T) {
	seen := make(map[string]struct{}, len(supportedComponents))

	for _, name := range SupportedComponentNames() {
		if name == "" {
			t.Fatal("supported component name should not be empty")
		}
		if _, exists := seen[name]; exists {
			t.Fatalf("duplicate supported component: %s", name)
		}
		seen[name] = struct{}{}

		if !IsSupportedComponent(name) {
			t.Fatalf("component should be registered: %s", name)
		}
	}
}

func TestAllComponentConstantsAreRegistered(t *testing.T) {
	expected := []string{
		ComponentInput,
		ComponentInputNumber,
		ComponentTextarea,
		ComponentSwitch,
		ComponentRadio,
		ComponentSelect,
		ComponentTreeSelectSingle,
		ComponentTreeSelectMulti,
		ComponentSelectMulti,
		ComponentImageUpload,
		ComponentFileUpload,
		ComponentRichText,
		ComponentJsonEditor,
		ComponentPassword,
		ComponentInputUrl,
		ComponentDateTimePicker,
		ComponentIconPicker,
	}

	got := SupportedComponentNames()
	if len(got) != len(expected) {
		t.Fatalf("supported component count mismatch: got=%d want=%d", len(got), len(expected))
	}

	for index := range expected {
		if got[index] != expected[index] {
			t.Fatalf("supported component mismatch at %d: got=%s want=%s", index, got[index], expected[index])
		}
	}
}
