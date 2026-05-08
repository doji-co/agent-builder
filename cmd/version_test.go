package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCommand_PrintsConfiguredVersion(t *testing.T) {
	SetVersion("v1.2.3")

	cmd := newVersionCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(buf.String(), "v1.2.3") {
		t.Fatalf("version output = %q, want to contain %q", buf.String(), "v1.2.3")
	}
}
