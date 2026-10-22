package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/xiaoxiaolin0918/rabbitmq-cli/internal/admin"
	"github.com/xiaoxiaolin0918/rabbitmq-cli/internal/output"
	"github.com/xiaoxiaolin0918/rabbitmq-cli/internal/risk"
	"github.com/xiaoxiaolin0918/rabbitmq-cli/internal/version"
)

type App struct {
	flagBin    string
	flagConfig string
	flagNode   string
	flagDryRun bool
	flagYes    bool
	flagFormat string
}

func Main() int {
	root := NewRoot()
	err := root.Execute()
	if err == nil {
		return 0
	}
	var ee *exitError
	if errors.As(err, &ee) {
		return ee.code
	}
	fmt.Fprintln(os.Stderr, "Error:", err)
	return 1
}

type exitError struct {
	msg  string
	code int
}

func (e *exitError) Error() string { return e.msg }

func NewRoot() *cobra.Command {
	a := &App{}
	root := &cobra.Command{
		Use:           "rabbitmq",
		Short:         "Thin RabbitMQ ops wrapper over rabbitmqadmin v2",
		Long:          "rabbitmq-cli thin-wraps official rabbitmqadmin. High-risk verbs need --dry-run or --yes (exit 10).",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version.Version,
	}
	root.SetVersionTemplate("{\"ok\":true,\"result\":{\"name\":\"rabbitmq-cli\",\"version\":\"{{.Version}}\"}}\n")
	root.CompletionOptions.DisableDefaultCmd = true

	p := root.PersistentFlags()
	p.StringVar(&a.flagBin, "bin", "", "path to rabbitmqadmin.exe")
	p.StringVarP(&a.flagConfig, "config", "c", "", "rabbitmqadmin config file (official -c)")
	p.StringVarP(&a.flagNode, "node", "N", "", "rabbitmqadmin node alias (official -N)")
	p.BoolVar(&a.flagDryRun, "dry-run", false, "print would-exec JSON; do not call rabbitmqadmin")
	p.BoolVarP(&a.flagYes, "yes", "y", false, "confirm a high-risk write")
	p.StringVar(&a.flagFormat, "format", "json", "wrapper output format for doctor/skills: json|pretty")

	root.AddCommand(a.newDoctorCmd(), a.newSkillsCmd(), a.newAdminCmd(), a.newCtlCmd())
	return root
}

func (a *App) prependConn(args []string) []string {
	out := []string{}
	if a.flagConfig != "" {
		out = append(out, "--config", a.flagConfig)
	}
	if a.flagNode != "" {
		out = append(out, "--node", a.flagNode)
	}
	return append(out, args...)
}

func (a *App) runAdmin(args []string) error {
	bin, err := admin.FindRabbitmqadmin(a.flagBin)
	if err != nil {
		return err
	}
	argv := a.prependConn(args)
	if risk.HighRisk(args) {
		if a.flagDryRun {
			output.Would("rabbitmqadmin", risk.Redact(append([]string{bin}, argv...)))
			return nil
		}
		if !a.flagYes {
			output.OK(map[string]any{
				"blocked": true,
				"reason":  "high-risk write requires --dry-run or --yes",
				"argv":    risk.Redact(argv),
			})
			return &exitError{msg: "high-risk write requires --dry-run or --yes", code: 10}
		}
	}
	cmd := exec.Command(bin, argv...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return &exitError{msg: err.Error(), code: ee.ExitCode()}
		}
		return err
	}
	return nil
}
