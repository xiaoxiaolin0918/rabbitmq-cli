package output

import (
	"encoding/json"
	"fmt"
	"os"
)

func OK(result any) {
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"ok": true, "result": result})
}

func Fail(err error, code int) int {
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"ok": false, "error": err.Error()})
	return code
}

func Would(action string, argv []string) {
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "dry_run": true, "would": map[string]any{"action": action, "argv": argv},
	})
}

func Errf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
}
