package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gdmc",
	Short: "Go Daily Message Creator - Generate messages from git commits",
	Long: `Go Daily Message Creator (gdmc) is a global tool that generates personalized messages 
	based on git commits in any repository. 

	Simply navigate to any git repository and run gdmc commands to analyze commit history
	and generate professional reports or meeting transcripts using Gemini AI.

	Usage:
	cd /path/to/your/git/repository
	gdmc generate --author="your@email.com" --type=report --interval=daily`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(configCmd)
}
