package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var logN int

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show recent wiki activity",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		w, err := openWiki()
		if err != nil {
			return err
		}
		entries := readLogEntries(w)
		if len(entries) == 0 {
			fmt.Println("No activity yet.")
			return nil
		}
		start := len(entries) - logN
		if start < 0 {
			start = 0
		}
		for _, e := range entries[start:] {
			fmt.Println(e)
		}
		return nil
	},
}

func init() {
	logCmd.Flags().IntVarP(&logN, "number", "n", 20, "number of entries to show")
	rootCmd.AddCommand(logCmd)
}
