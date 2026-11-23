package cli

import (
	"math"
	"testing"
)

func TestParseHundredths(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input   string
		want    int64
		wantErr bool
	}{
		{input: "8", want: 800},
		{input: "8.2", want: 820},
		{input: "8.25", want: 825},
		{input: "-.01", want: -1},
		{input: "+0", want: 0},
		{input: "92233720368547758.07", want: math.MaxInt64},
		{input: "-92233720368547758.08", want: math.MinInt64},
		{input: "92233720368547758.08", wantErr: true},
		{input: "8.259", wantErr: true},
		{input: "8.", wantErr: true},
		{input: "8,25", wantErr: true},
		{input: "1e2", wantErr: true},
		{input: "+", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()

			got, err := parseHundredths(tc.input)

			if got != tc.want || (err != nil) != tc.wantErr {
				t.Errorf("parseHundredths(%q) = (%d, %v), want %d with error = %t", tc.input, got, err, tc.want, tc.wantErr)
			}
		})
	}
}
