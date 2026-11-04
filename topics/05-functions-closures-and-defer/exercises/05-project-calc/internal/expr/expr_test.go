package expr

import (
	"errors"
	"strings"
	"testing"
)

func TestEval(t *testing.T) {
	t.Parallel()
	tests := []struct {
		src  string
		want int64
	}{
		{"42", 42},
		{"  42  ", 42},
		{"1+2*3", 7},
		{"(1+2)*3", 9},
		{"2 * (3 + 4) * 5", 70},
		{"10-4-3", 3},
		{"100/10/5", 2},
		{"7/2", 3},
		{"-7/2", -3},
		{"7%3", 1},
		{"-7%3", -1},
		{"-3", -3},
		{"--3", 3},
		{"-(2+3)*2", -10},
		{"2 * (3 + 4) - -1", 15},
		{"((((7))))", 7},
		{"9223372036854775807", 9223372036854775807},
	}
	for _, tt := range tests {
		t.Run(tt.src, func(t *testing.T) {
			t.Parallel()
			got, err := Eval(tt.src)
			if err != nil {
				t.Fatalf("Eval(%q) error = %v, want %d", tt.src, err, tt.want)
			}
			if got != tt.want {
				t.Errorf("Eval(%q) = %d, want %d", tt.src, got, tt.want)
			}
		})
	}
}

func TestEvalErrors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		src     string
		wantErr string
	}{
		{"", "unexpected end of input"},
		{"1 +", "unexpected end of input"},
		{"(1+2", "missing ')' at offset 4"},
		{"(1+2 3", "missing ')' at offset 5"},
		{"1+2)", "unexpected ')' at offset 3"},
		{"1 2", "unexpected 2 at offset 2"},
		{"2**3", "unexpected '*' at offset 2"},
		{"3 $ 4", "unexpected character '$' at offset 2"},
		{"99999999999999999999", `number "99999999999999999999" at offset 0 is out of range`},
		{"4/0", "offset 1: division by zero"},
		{"4 % (2-2)", "offset 2: division by zero"},
	}
	for _, tt := range tests {
		t.Run(tt.src, func(t *testing.T) {
			t.Parallel()
			got, err := Eval(tt.src)
			if err == nil {
				t.Fatalf("Eval(%q) = %d, nil; want error %q", tt.src, got, tt.wantErr)
			}
			if err.Error() != tt.wantErr {
				t.Errorf("Eval(%q) error = %q, want %q", tt.src, err, tt.wantErr)
			}
			if got != 0 {
				t.Errorf("Eval(%q) = %d with an error, want 0", tt.src, got)
			}
		})
	}
}

func TestEvalDivisionByZeroIs(t *testing.T) {
	t.Parallel()
	for _, src := range []string{"1/0", "5 % 0", "(1+1) / (3-3)"} {
		if _, err := Eval(src); !errors.Is(err, ErrDivisionByZero) {
			t.Errorf("Eval(%q) error = %v, want one wrapping ErrDivisionByZero", src, err)
		}
	}
}

func TestTrace(t *testing.T) {
	t.Parallel()
	tests := []struct {
		src       string
		want      string
		wantValue int64
		wantErr   bool
	}{
		{
			src: "1+2",
			want: `> expr
  > term
    > factor
    < factor = 1
  < term = 1
  > term
    > factor
    < factor = 2
  < term = 2
< expr = 3
`,
			wantValue: 3,
		},
		{
			src: "-2*3",
			want: `> expr
  > term
    > factor
      > factor
      < factor = 2
    < factor = -2
    > factor
    < factor = 3
  < term = -6
< expr = -6
`,
			wantValue: -6,
		},
		{
			src: "(1",
			want: `> expr
  > term
    > factor
      > expr
        > term
          > factor
          < factor = 1
        < term = 1
      < expr = 1
    < factor error
  < term error
< expr error
`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.src, func(t *testing.T) {
			t.Parallel()
			var sb strings.Builder
			v, err := Trace(tt.src, &sb)
			if (err != nil) != tt.wantErr || v != tt.wantValue {
				t.Errorf("Trace(%q) = %d, %v; want %d, error: %t", tt.src, v, err, tt.wantValue, tt.wantErr)
			}
			if got := sb.String(); got != tt.want {
				t.Errorf("Trace(%q) wrote:\n%s\nwant:\n%s", tt.src, got, tt.want)
			}
		})
	}
}
