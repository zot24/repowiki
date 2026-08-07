package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zot24/repowiki/internal/wiki"
)

var skillCmd = &cobra.Command{
	Use:   "skill [install]",
	Short: "Print the agent skill, or `skill install` it to .claude/skills/repowiki/",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Print(wiki.SkillTemplate)
			return nil
		}
		if args[0] != "install" {
			return fmt.Errorf("unknown subcommand %q (want `repowiki skill` or `repowiki skill install`)", args[0])
		}
		w, err := openWiki()
		if err != nil {
			return err
		}
		result, err := wiki.InstallSkill(w.RepoRoot)
		if err != nil {
			return err
		}
		fmt.Printf(".claude/skills/repowiki/SKILL.md: %s\n", result)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(skillCmd)
}
