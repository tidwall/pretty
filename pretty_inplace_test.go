package pretty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestUglyInPlaceEscapedStrings(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"backslash", `["\\", "a b"]`},
		{"backslashes", `["\\\\", "a b"]`},
		{"escaped_quote", `["a\" b", "c d"]`},
		{"backslash_and_quote", `["\\\" x", "a b"]`},
		{"object", `{"a":"\\", "b":"c d"}`},
		{"nested", `{"\\":"a b", "nested": [{"quote":"\\\"", "text":"c d"}]}`},
		{"other_escapes", `["\u005c", "a b", "\t\n\r\b\f\/"]`},
		{"plain", `["a b", true, null, 12.30]`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", " ", "  ", "   ", "        ", "\t\n  "} {
				t.Run(fmt.Sprintf("leading_whitespace_%d", len(prefix)), func(t *testing.T) {
					input := []byte(prefix + tt.input + " \n")
					var want bytes.Buffer
					if err := json.Compact(&want, input); err != nil {
						t.Fatal(err)
					}
					if got := Ugly(input); !bytes.Equal(got, want.Bytes()) {
						t.Fatalf("Ugly: got %q, want %q", got, want.Bytes())
					}
					got := UglyInPlace(input)
					if !bytes.Equal(got, want.Bytes()) {
						t.Fatalf("UglyInPlace: got %q, want %q", got, want.Bytes())
					}
					if &got[0] != &input[0] {
						t.Fatal("UglyInPlace did not reuse the input buffer")
					}
				})
			}
		})
	}
}
