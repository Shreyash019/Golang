package main

import (
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"no coupon", []string{"1999", "500"}, "items:    2\nsubtotal: 2499\ndiscount: 0\ntotal:    2499\n"},
		{"SAVE10 rounds down", []string{"-coupon", "SAVE10", "1999", "500"}, "items:    2\nsubtotal: 2499\ndiscount: 249\ntotal:    2250\n"},
		{"FLAT500", []string{"-coupon", "FLAT500", "1200"}, "items:    1\nsubtotal: 1200\ndiscount: 500\ntotal:    700\n"},
		{"FLAT500 never below zero", []string{"-coupon=FLAT500", "300"}, "items:    1\nsubtotal: 300\ndiscount: 300\ntotal:    0\n"},
		{"BOGO frees the cheapest", []string{"-coupon", "BOGO", "800", "250", "990"}, "items:    3\nsubtotal: 2040\ndiscount: 250\ntotal:    1790\n"},
		{"BOGO needs two items", []string{"-coupon", "BOGO", "800"}, "items:    1\nsubtotal: 800\ndiscount: 0\ntotal:    800\n"},
		{"free items are fine", []string{"0", "0"}, "items:    2\nsubtotal: 0\ndiscount: 0\ntotal:    0\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var sb strings.Builder
			if err := run(tt.args, &sb); err != nil {
				t.Fatalf("run(%q) error = %v", tt.args, err)
			}
			if got := sb.String(); got != tt.want {
				t.Errorf("run(%q) printed:\n%s\nwant:\n%s", tt.args, got, tt.want)
			}
		})
	}
}

func TestRunErrors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		args    []string
		wantErr string
	}{
		{nil, "no items"},
		{[]string{"-coupon", "SAVE10"}, "no items"},
		{[]string{"12.50"}, "invalid price"},
		{[]string{"-5"}, ""}, // parsed as an unknown flag: any error will do
		{[]string{"100", "-5"}, "invalid price"},
		{[]string{"abc"}, "invalid price"},
		{[]string{"-coupon", "HALFOFF", "100"}, "unknown coupon"},
	}
	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			t.Parallel()
			var sb strings.Builder
			err := run(tt.args, &sb)
			if err == nil {
				t.Fatalf("run(%q) = nil, printed %q; want an error", tt.args, sb.String())
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("run(%q) error = %q, want it to mention %q", tt.args, err, tt.wantErr)
			}
		})
	}
}
