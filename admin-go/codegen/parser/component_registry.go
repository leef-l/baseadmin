package parser

var supportedComponents = []string{
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

var supportedComponentSet = func() map[string]struct{} {
	set := make(map[string]struct{}, len(supportedComponents))
	for _, name := range supportedComponents {
		set[name] = struct{}{}
	}
	return set
}()

func IsSupportedComponent(name string) bool {
	_, ok := supportedComponentSet[name]
	return ok
}

func SupportedComponentNames() []string {
	return append([]string(nil), supportedComponents...)
}
