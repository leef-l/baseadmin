package password

import "testing"

func TestValidatePolicy(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr string
	}{
		{name: "blank", value: "   ", wantErr: "密码不能为空"},
		{name: "has spaces", value: "abc 12345", wantErr: "密码不能包含空白字符"},
		{name: "too short", value: "ab12", wantErr: "密码长度需为8-64位"},
		{name: "letters only", value: "abcdefgh", wantErr: "密码必须同时包含字母和数字"},
		{name: "digits only", value: "12345678", wantErr: "密码必须同时包含字母和数字"},
		{name: "valid", value: "abc12345", wantErr: ""},
	}

	for _, tc := range tests {
		err := ValidatePolicy(tc.value)
		if tc.wantErr == "" {
			if err != nil {
				t.Fatalf("%s: unexpected error: %v", tc.name, err)
			}
			continue
		}
		if err == nil || err.Error() != tc.wantErr {
			t.Fatalf("%s: got err=%v want=%q", tc.name, err, tc.wantErr)
		}
	}
}
