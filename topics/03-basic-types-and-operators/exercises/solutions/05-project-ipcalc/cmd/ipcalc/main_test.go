package main

import (
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		args []string
		want string
	}{
		{
			args: []string{"192.168.1.10/24"},
			want: `address:   192.168.1.10
network:   192.168.1.0/24
netmask:   255.255.255.0
broadcast: 192.168.1.255
hosts:     254 (192.168.1.1 - 192.168.1.254)
`,
		},
		{
			args: []string{"10.0.0.0/8", "10.20.30.40"},
			want: `address:   10.0.0.0
network:   10.0.0.0/8
netmask:   255.0.0.0
broadcast: 10.255.255.255
hosts:     16777214 (10.0.0.1 - 10.255.255.254)
10.20.30.40 is in 10.0.0.0/8
`,
		},
		{
			args: []string{"172.16.5.4/31", "172.16.5.9"},
			want: `address:   172.16.5.4
network:   172.16.5.4/31
netmask:   255.255.255.254
broadcast: 172.16.5.5
hosts:     2 (172.16.5.4 - 172.16.5.5)
172.16.5.9 is not in 172.16.5.4/31
`,
		},
	}
	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			var sb strings.Builder
			if err := run(&sb, tt.args); err != nil {
				t.Fatalf("run(%q) returned error %v", tt.args, err)
			}
			if got := sb.String(); got != tt.want {
				t.Errorf("run(%q) printed:\n%s\nwant:\n%s", tt.args, got, tt.want)
			}
		})
	}
}

func TestRunErrors(t *testing.T) {
	for _, args := range [][]string{
		nil,
		{"10.0.0.0/8", "10.0.0.1", "extra"},
		{"10.0.0.0"},
		{"10.0.0.0/8", "10.0.0"},
	} {
		var sb strings.Builder
		if err := run(&sb, args); err == nil {
			t.Errorf("run(%q) = nil, printed %q; want an error", args, sb.String())
		}
	}
}
