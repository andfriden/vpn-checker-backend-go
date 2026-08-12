package validator

import "testing"

func TestValidate(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid vless",
			input: "vless://uuid@example.com:443",
			want:  true,
		},
		{
			name:  "valid trojan",
			input: "trojan://password@example.com:443",
			want:  true,
		},
		{
			name:  "invalid text",
			input: "hello world",
			want:  false,
		},
		{
			name:  "invalid scheme",
			input: "ftp://example.com",
			want:  false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			got := Validate(tt.input)

			if got != tt.want {
				t.Fatalf(
					"expected %v got %v",
					tt.want,
					got,
				)
			}

		})
	}
}
