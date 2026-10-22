package admin

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFindRabbitmqadminExplicitMissing(t *testing.T) {
	_, err := FindRabbitmqadmin(filepath.Join(t.TempDir(), "nope.exe"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFindRabbitmqadminDefaultWindowsPath(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows path probe")
	}
	p := `D:\tools\rabbitmqadmin\rabbitmqadmin.exe`
	if _, err := os.Stat(p); err != nil {
		t.Skip("official binary not installed on this machine")
	}
	got, err := FindRabbitmqadmin("")
	if err != nil {
		t.Fatal(err)
	}
	if got != p && filepath.Base(got) != "rabbitmqadmin.exe" {
		t.Fatalf("unexpected bin %s", got)
	}
}
