package risk

import "strings"

// globalValueFlags are rabbitmqadmin v2 global options that consume a value.
// They may appear before the command group (and, being clap globals, between
// group and verb), so they must be skipped when locating the command words.
var globalValueFlags = map[string]bool{
	"-c": true, "--config": true,
	"-N": true, "--node": true,
	"-H": true, "--host": true,
	"-P": true, "--port": true,
	"-U": true, "--base-uri": true,
	"--path-prefix": true,
	"-V":            true, "--vhost": true,
	"-u": true, "--username": true,
	"-p": true, "--password": true,
	"--tls-ca-cert-file": true, "--tls-cert-file": true, "--tls-key-file": true,
	"--timeout":     true,
	"--table-style": true,
}

// commandWords returns up to max positional command words, skipping flags
// (and the values of known global value flags).
func commandWords(args []string, max int) []string {
	words := []string{}
	for i := 0; i < len(args) && len(words) < max; i++ {
		a := args[i]
		if a == "--" {
			continue
		}
		if strings.HasPrefix(a, "-") && len(a) > 1 {
			if !strings.Contains(a, "=") && globalValueFlags[a] {
				i++
			}
			continue
		}
		words = append(words, strings.ToLower(a))
	}
	return words
}

var topLevelWrite = map[string]bool{
	"delete": true, "purge": true, "close": true, "declare": true,
	"import": true, "export": true, "publish": true, "rebalance": true,
	"get": true, // polling get can consume/ack messages
}

var subWrite = map[string]bool{
	"declare": true, "delete": true, "purge": true, "close": true,
	"publish": true, "import": true, "bind": true, "unbind": true,
	"set": true, "update": true, "clear": true, "patch": true,
	"enable": true, "enable_all": true, "rebalance": true, "queues": true,
}

var subWritePrefixes = []string{
	"declare_", "delete_", "update_", "close_", "import_",
	"enable_", "disable_", "add_",
}

func isWriteVerb(sub string) bool {
	if subWrite[sub] {
		return true
	}
	for _, p := range subWritePrefixes {
		if strings.HasPrefix(sub, p) {
			return true
		}
	}
	return false
}

// HighRisk reports whether rabbitmqadmin argv looks like a write / destructive call.
// Leading global flags (e.g. --vhost /, --non-interactive) are skipped.
func HighRisk(args []string) bool {
	w := commandWords(args, 3)
	if len(w) == 0 {
		return false
	}
	head := w[0]
	if topLevelWrite[head] {
		return true
	}
	if len(w) < 2 {
		return false
	}
	sub := w[1]
	switch head {
	case "list", "show", "health_check", "help", "channels", "nodes", "plugins",
		"deprecated_features", "auth_attempts", "passwords", "shell":
		return false
	case "definitions":
		// export/export_from_vhost write a local file containing password hashes.
		return strings.HasPrefix(sub, "export") || strings.HasPrefix(sub, "import")
	case "rebalance":
		return true
	case "tanzu":
		return len(w) > 2 && isWriteVerb(w[2])
	}
	if sub == "queues" { // only meaningful as `rebalance queues`
		return false
	}
	return isWriteVerb(sub)
}

// Redact masks secret flag values (passwords) in argv before printing.
func Redact(argv []string) []string {
	out := make([]string, len(argv))
	copy(out, argv)
	for i := 0; i < len(out); i++ {
		a := out[i]
		switch {
		case a == "--password" || a == "-p":
			if i+1 < len(out) {
				out[i+1] = "***"
				i++
			}
		case strings.HasPrefix(a, "--password="):
			out[i] = "--password=***"
		}
	}
	return out
}
