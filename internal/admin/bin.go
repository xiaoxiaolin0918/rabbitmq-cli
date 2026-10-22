package admin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func FindRabbitmqadmin(explicit string) (string, error) {
	if explicit != "" {
		if st, err := os.Stat(explicit); err == nil && !st.IsDir() {
			return explicit, nil
		}
		return "", fmt.Errorf("rabbitmqadmin not found at %s", explicit)
	}
	if env := os.Getenv("RABBITMQ_CLI_ADMIN_BIN"); env != "" {
		if st, err := os.Stat(env); err == nil && !st.IsDir() {
			return env, nil
		}
	}
	candidates := []string{}
	if runtime.GOOS == "windows" {
		candidates = append(candidates,
			`D:\tools\rabbitmqadmin\rabbitmqadmin.exe`,
			filepath.Join(os.Getenv("LOCALAPPDATA"), "rabbitmqadmin", "rabbitmqadmin.exe"),
		)
	}
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			return c, nil
		}
	}
	if p, err := exec.LookPath("rabbitmqadmin"); err == nil {
		return p, nil
	}
	if p, err := exec.LookPath("rabbitmqadmin.exe"); err == nil {
		return p, nil
	}
	return "", fmt.Errorf("rabbitmqadmin not found; set RABBITMQ_CLI_ADMIN_BIN or install to D:\\tools\\rabbitmqadmin")
}

func Version(bin string) (string, error) {
	out, err := exec.Command(bin, "--version").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}
