package cmd

import (
	"strings"
	"testing"
)

func TestDefaultSkillsDirHasNoControlChars(t *testing.T) {
	if strings.ContainsAny(defaultSkillsDir, "\r\n\t") {
		t.Fatalf("defaultSkillsDir contains control chars: %q", defaultSkillsDir)
	}
	if !strings.HasSuffix(strings.ReplaceAll(defaultSkillsDir, `\`, "/"), "work/ops/rabbitmq-cli/skills") {
		t.Fatalf("unexpected defaultSkillsDir %q", defaultSkillsDir)
	}
}

func TestValidSkillName(t *testing.T) {
	for _, ok := range []string{"rabbitmq-shared", "rabbitmq-admin-ops"} {
		if !validSkillName(ok) {
			t.Errorf("%q should be valid", ok)
		}
	}
	for _, bad := range []string{"", ".", "..", `..\x`, "../x", `a\b`, "a/b", `C:x`} {
		if validSkillName(bad) {
			t.Errorf("%q should be invalid", bad)
		}
	}
}

func TestStripWrapperFlags(t *testing.T) {
	a := &App{}
	got := stripWrapperFlags(a, []string{"-N", "lab", "--dry-run", "queues", "list"})
	if strings.Join(got, " ") != "queues list" || a.flagNode != "lab" || !a.flagDryRun {
		t.Fatalf("got %v node=%q dry=%v", got, a.flagNode, a.flagDryRun)
	}
}
