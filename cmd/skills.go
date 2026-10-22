package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/xiaoxiaolin0918/rabbitmq-cli/internal/output"
)

func (a *App) newSkillsCmd() *cobra.Command {
	c := &cobra.Command{Use: "skills", Short: "List or read bundled Trae-style skills"}
	c.AddCommand(
		&cobra.Command{Use: "list", Short: "List skills", RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := skillsDir()
			if err != nil {
				return err
			}
			entries, err := os.ReadDir(dir)
			if err != nil {
				return err
			}
			names := []string{}
			for _, e := range entries {
				if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
					names = append(names, e.Name())
				}
			}
			output.OK(map[string]any{"dir": dir, "skills": names})
			return nil
		}},
		&cobra.Command{Use: "path", Short: "Print skills directory", RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := skillsDir()
			if err != nil {
				return err
			}
			output.OK(map[string]any{"path": dir})
			return nil
		}},
		&cobra.Command{Use: "read [name]", Short: "Print a skill markdown", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := skillsDir()
			if err != nil {
				return err
			}
			if !validSkillName(args[0]) {
				return fmt.Errorf("invalid skill name %q", args[0])
			}
			p := filepath.Join(dir, args[0], "SKILL.md")
			b, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			fmt.Print(string(b))
			return nil
		}},
	)
	return c
}

// defaultSkillsDir is the checkout fallback when the exe is copied elsewhere.
var defaultSkillsDir = filepath.Join(`D:\`, "work", "ops", "rabbitmq-cli", "skills")

// validSkillName rejects path separators / traversal in `skills read <name>`.
func validSkillName(name string) bool {
	return name != "" && name != "." && name != ".." && !strings.ContainsAny(name, `/\:`) && !strings.Contains(name, "..")
}

func skillsDir() (string, error) {
	candidates := []string{}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "skills"))
	}
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(wd, "skills"))
	}
	candidates = append(candidates, defaultSkillsDir)
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && st.IsDir() {
			return c, nil
		}
	}
	return "", fmt.Errorf("skills directory not found")
}
