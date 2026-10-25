package report

import (
	"errors"
	"strings"
	"testing"

	"github.com/Shreyash019/Golang/topics/01-go-fundamentals-and-toolchain/exercises/solutions/05-project-pkgscan/internal/scan"
)

func TestWrite(t *testing.T) {
	tests := []struct {
		name string
		pkgs []scan.Package
		want string
	}{
		{
			name: "no packages",
			pkgs: nil,
			want: "DIR  PACKAGE  KIND  FILES  TESTS  IMPORTS  DOC\n\n0 packages (0 commands, 0 libraries)\n",
		},
		{
			name: "command and library",
			pkgs: []scan.Package{
				{Dir: "cmd/api", Name: "main", Files: 1, Imports: []string{"example.com/app/internal/order", "fmt"}},
				{Dir: "internal/order", Name: "order", Files: 2, TestFiles: 1, Imports: []string{"fmt", "sync"}, HasDoc: true},
			},
			want: "" +
				"DIR             PACKAGE  KIND     FILES  TESTS  IMPORTS  DOC\n" +
				"cmd/api         main     command  1      0      2        no\n" +
				"internal/order  order    library  2      1      2        yes\n" +
				"\n" +
				"2 packages (1 command, 1 library)\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			if err := Write(&sb, tt.pkgs); err != nil {
				t.Fatalf("Write() unexpected error: %v", err)
			}
			if got := sb.String(); got != tt.want {
				t.Errorf("Write() output:\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}

type failWriter struct{}

func (failWriter) Write([]byte) (int, error) { return 0, errors.New("disk full") }

func TestWriteReturnsWriteErrors(t *testing.T) {
	err := Write(failWriter{}, []scan.Package{{Dir: ".", Name: "x", Files: 1}})
	if err == nil || !strings.Contains(err.Error(), "disk full") {
		t.Errorf("Write(failing writer) error = %v, want it to contain %q", err, "disk full")
	}
}
