package service

import "testing"

func TestFormatWithCommas(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want string
	}{
		{name: "zero", in: 0, want: "0"},
		{name: "single digit", in: 7, want: "7"},
		{name: "three digits", in: 999, want: "999"},
		{name: "thousands", in: 1234, want: "1,234"},
		{name: "tens of thousands", in: 12345, want: "12,345"},
		{name: "millions", in: 1234567, want: "1,234,567"},
		{name: "negative", in: -1234567, want: "-1,234,567"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := formatWithCommas(tc.in)
			if got != tc.want {
				t.Fatalf("formatWithCommas(%d) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
