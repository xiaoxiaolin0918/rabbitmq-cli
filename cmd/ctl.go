package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/xiaoxiaolin0918/rabbitmq-cli/internal/output"
)

func (a *App) newCtlCmd() *cobra.Command {
	return &cobra.Command{
		Use:                "ctl",
		Short:              "Passthrough to rabbitmqctl when installed on PATH",
		DisableFlagParsing: true,
		SilenceUsage:       true,
		RunE: func(cmd *cobra.Command, args []string) error {
			bin, err := exec.LookPath("rabbitmqctl")
			if err != nil {
				bin, err = exec.LookPath("rabbitmqctl.bat")
			}
			if err != nil {
				output.OK(map[string]any{
					"bin_ok": false,
					"hint":   "rabbitmqctl not on PATH; install RabbitMQ server tools or skip ctl",
				})
				return &exitError{msg: "rabbitmqctl not found", code: 1}
			}
			if len(args) == 0 {
				args = []string{"help"}
			}
			c := exec.Command(bin, args...)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			if err := c.Run(); err != nil {
				var ee *exec.ExitError
				if errors.As(err, &ee) {
					return &exitError{msg: err.Error(), code: ee.ExitCode()}
				}
				return err
			}
			return nil
		},
	}
}
