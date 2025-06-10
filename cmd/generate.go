package cmd

import (
	"fmt"
	"strings"

	"gdmc/internal/config"
	"gdmc/internal/gemini"
	"gdmc/internal/git"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate message from git commits",
	Long:  "Generate a message using Gemini AI based on your git commits for a specified time period",
	RunE:  runGenerate,
}

var (
	author      string
	interval    string
	template    string
	apiKey      string
	language    string
)

func init() {
	generateCmd.Flags().StringVarP(&author, "author", "a", "", "Git author email")
	generateCmd.Flags().StringVarP(&interval, "interval", "i", "", "Time interval (daily, weekly, monthly)")
	generateCmd.Flags().StringVarP(&template, "template", "t", "", "Message template (report, transcript, summary)")
	generateCmd.Flags().StringVar(&apiKey, "api-key", "", "Gemini API key")
	generateCmd.Flags().StringVarP(&language, "language", "l", "", "Output language (en, tr)")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if author == "" {
		author = cfg.Author
		if author == "" {
			return fmt.Errorf("author email is required. Use --author flag or set in config")
		}
	}

	if template == "" {
		template = cfg.DefaultType
	}

	if interval == "" {
		interval = "daily"
	}

	if language == "" {
		language = cfg.Language
	}

	if apiKey == "" {
		apiKey = cfg.APIKey
		if apiKey == "" {
			return fmt.Errorf("gemini API key is required. Use --api-key flag or set GEMINI_API_KEY env var")
		}
	}

	if _, exists := cfg.Templates[template]; !exists {
		return fmt.Errorf("unknown template: %s. Available: %v", template, getTemplateKeys(cfg.Templates))
	}

	if _, exists := cfg.Intervals[interval]; !exists {
		return fmt.Errorf("unknown interval: %s. Available: %v", interval, getIntervalKeys(cfg.Intervals))
	}

	commitFetcher := git.NewCommitFetcher(cfg)
	commits, err := commitFetcher.FetchCommits(author, interval)
	if err != nil {
		return fmt.Errorf("failed to fetch commits: %w", err)
	}

	if len(commits) == 0 {
		fmt.Printf("No commits found for %s in %s period.\n", author, interval)
		return nil
	}

	fmt.Printf("Found %d commits for %s period\n", len(commits), interval)
	fmt.Printf("Language: %s\n", language)

	geminiClient := gemini.NewClient(apiKey, cfg)
	message, err := geminiClient.GenerateMessage(commits, template, interval, author, language)
	if err != nil {
		return fmt.Errorf("failed to generate message: %w", err)
	}

	fmt.Printf("\n%s (%s)\n", cfg.Templates[template].Name, cfg.Intervals[interval].Name)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println(message)

	return nil
}

func getTemplateKeys(m map[string]config.Template) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func getIntervalKeys(m map[string]config.Interval) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
