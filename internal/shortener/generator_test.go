package shortener

import "testing"

func TestEncode(t *testing.T) {
	gen := NewGenerator(0)

	tests := []struct {
		id   uint64
		want string
	}{
		{1, "b"},
		{61, "9"},
		{62, "ba"},
		{3844, "baa"},
	}

	for _, tt := range tests {
		if got := gen.Encode(tt.id); got != tt.want {
			t.Errorf("Encode(%d) = %q, want %q", tt.id, got, tt.want)
		}
	}
}

func TestGenerateUnique(t *testing.T) {
	gen := NewGenerator(0)
	seen := make(map[string]bool)

	for range 1000 {
		code := gen.Generate()
		if seen[code] {
			t.Fatalf("duplicate short code: %s", code)
		}
		seen[code] = true
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		url  string
		want bool
	}{
		{"https://example.com", true},
		{"http://example.com/path", true},
		{"ftp://example.com", false},
		{"not-a-url", false},
		{"", false},
	}

	for _, tt := range tests {
		err := validateURL(tt.url)
		got := err == nil
		if got != tt.want {
			t.Errorf("validateURL(%q) = %v, want %v", tt.url, got, tt.want)
		}
	}
}
