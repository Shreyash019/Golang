package main

import (
	"strings"
	"testing"
)

func TestPrediction(t *testing.T) {
	if prediction == "" {
		t.Fatal("prediction.go is empty: write your prediction before running anything")
	}
	var sb strings.Builder
	run(&sb)
	if got := sb.String(); got != prediction {
		t.Errorf("your prediction differs from the real output.\n--- real output ---\n%s--- your prediction ---\n%s\nExplain every line that differs (notes §5.3–5.8 and §7.1).", got, prediction)
	}
}
