package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestPatternsCommand_PrintsAllPatterns(t *testing.T) {
	cmd := newPatternsCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	expectedStrings := []string{
		"Sequential",
		"Parallel",
		"LLM-Coordinated",
		"Loop",
	}
	for _, expected := range expectedStrings {
		if !strings.Contains(buf.String(), expected) {
			t.Fatalf("patterns output = %q, want to contain %q", buf.String(), expected)
		}
	}
}
