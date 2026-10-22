package cmd

import (
	"github.com/spf13/cobra"
)

func (a *App) newAdminCmd() *cobra.Command {
	c := &cobra.Command{
		Use:                "admin",
		Short:              "Passthrough to official rabbitmqadmin",
		DisableFlagParsing: true,
		SilenceUsage:       true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// When DisableFlagParsing is set, cobra still strips the command name;
			// global flags must be placed before `admin` by the user, or we accept
			// leading wrapper flags here.
			args = stripWrapperFlags(a, args)
			if len(args) == 0 {
				args = []string{"help"}
			}
			return a.runAdmin(args)
		},
	}
	return c
}

func stripWrapperFlags(a *App, args []string) []string {
	out := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--dry-run":
			a.flagDryRun = true
		case "--yes", "-y":
			a.flagYes = true
		case "--bin":
			if i+1 < len(args) {
				i++
				a.flagBin = args[i]
			}
		case "--config", "-c":
			if i+1 < len(args) {
				i++
				a.flagConfig = args[i]
			}
		case "--node", "-N":
			if i+1 < len(args) {
				i++
				a.flagNode = args[i]
			}
		case "--":
			out = append(out, args[i+1:]...)
			return out
		default:
			out = append(out, args[i])
		}
	}
	return out
}
