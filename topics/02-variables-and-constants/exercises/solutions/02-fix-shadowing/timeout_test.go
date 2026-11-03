package shadowing

import (
	"testing"
	"time"
)

func TestParseTimeout(t *testing.T) {
	t.Parallel()
	const def = 30 * time.Second
	tests := []struct {
		in      string
		want    time.Duration
		wantErr bool
	}{
		{in: "", want: def},
		{in: "5s", want: 5 * time.Second},
		{in: "1.5s", want: 1500 * time.Millisecond},
		{in: "250ms", want: 250 * time.Millisecond},
		{in: "1h30m", want: 90 * time.Minute},
		{in: "soon", wantErr: true},
		{in: "0s", wantErr: true},
		{in: "-1s", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			got, err := ParseTimeout(tt.in, def)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseTimeout(%q, %v) = %v, nil; want an error", tt.in, def, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseTimeout(%q, %v) returned error %v; want %v", tt.in, def, err, tt.want)
			}
			if got != tt.want {
				t.Errorf("ParseTimeout(%q, %v) = %v, want %v", tt.in, def, got, tt.want)
			}
		})
	}
}
