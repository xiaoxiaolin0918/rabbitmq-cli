package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/xiaoxiaolin0918/rabbitmq-cli/internal/admin"
	"github.com/xiaoxiaolin0918/rabbitmq-cli/internal/output"
	"github.com/xiaoxiaolin0918/rabbitmq-cli/internal/version"
)

func (a *App) newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Locate rabbitmqadmin and print environment hints",
		RunE: func(cmd *cobra.Command, args []string) error {
			bin, err := admin.FindRabbitmqadmin(a.flagBin)
			res := map[string]any{
				"wrapper_version": version.Version,
				"bin":             bin,
				"bin_ok":          err == nil,
			}
			if err != nil {
				res["bin_error"] = err.Error()
			} else if ver, e := admin.Version(bin); e == nil {
				res["admin_version"] = ver
			} else {
				res["admin_version_error"] = e.Error()
			}
			home, _ := os.UserHomeDir()
			cfg := filepath.Join(home, ".rabbitmqadmin.conf")
			res["default_config"] = cfg
			if st, e := os.Stat(cfg); e == nil && !st.IsDir() {
				res["config_exists"] = true
			} else {
				res["config_exists"] = false
				res["config_hint"] = "copy rabbitmqadmin.conf.example to %USERPROFILE%\\.rabbitmqadmin.conf"
			}
			if a.flagConfig != "" {
				res["config_flag"] = a.flagConfig
			}
			if a.flagNode != "" {
				res["node_flag"] = a.flagNode
			}
			if p, e := exec.LookPath("rabbitmqctl"); e == nil {
				res["ctl"] = p
				res["ctl_ok"] = true
			} else if p, e := exec.LookPath("rabbitmqctl.bat"); e == nil {
				res["ctl"] = p
				res["ctl_ok"] = true
			} else {
				res["ctl_ok"] = false
			}
			output.OK(res)
			if err != nil {
				return &exitError{msg: err.Error(), code: 1}
			}
			return nil
		},
	}
}
