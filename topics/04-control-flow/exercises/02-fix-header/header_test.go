package header

import (
	"slices"
	"testing"
)

func TestParse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		lines []string
		want  []Field
	}{
		{
			name:  "header only",
			lines: []string{"title = Weekly report", "owner=ops"},
			want:  []Field{{"title", "Weekly report"}, {"owner", "ops"}},
		},
		{
			name:  "comments, blanks and lines without =",
			lines: []string{"# generated", "", "  title = Q3  ", "not a field", "owner = ops"},
			want:  []Field{{"title", "Q3"}, {"owner", "ops"}},
		},
		{
			name:  "body after the separator is ignored",
			lines: []string{"title = Weekly report", "---", "status = draft", "x = 1"},
			want:  []Field{{"title", "Weekly report"}},
		},
		{
			name:  "separator with surrounding spaces",
			lines: []string{"a = 1", "  ---  ", "b = 2"},
			want:  []Field{{"a", "1"}},
		},
		{
			name:  "empty header",
			lines: []string{"---", "a = 1"},
			want:  nil,
		},
		{
			name:  "no lines",
			lines: nil,
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := Parse(tt.lines); !slices.Equal(got, tt.want) {
				t.Errorf("Parse(%q) = %v, want %v", tt.lines, got, tt.want)
			}
		})
	}
}
