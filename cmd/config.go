package cmd

import (
	"fmt"

	"gdmc/internal/config"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "View and manage dmc configuration settings",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE:  showConfig,
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set configuration value",
	Args:  cobra.ExactArgs(2),
	RunE:  setConfig,
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
}

func showConfig(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Printf("Author: %s\n", cfg.Author)
	fmt.Printf("Default Template: %s\n", cfg.DefaultType)
	fmt.Printf("API Key: %s\n", maskAPIKey(cfg.APIKey))

	fmt.Println("\nAvailable Intervals:")
	for key, interval := range cfg.Intervals {
		fmt.Printf("  %s: %s (%s to %s)\n", key, interval.Name, interval.Since, interval.Until)
	}

	fmt.Println("\nAvailable Templates:")
	for key, tmpl := range cfg.Templates {
		fmt.Printf("  %s: %s - %s\n", key, tmpl.Name, tmpl.Description)
	}

	return nil
}

func setConfig(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	key, value := args[0], args[1]

	switch key {
	case "author":
		cfg.Author = value
	case "default_type":
		if _, exists := cfg.Templates[value]; !exists {
			return fmt.Errorf("unknown template: %s", value)
		}
		cfg.DefaultType = value
	case "api_key":
		cfg.APIKey = value
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Set %s = %s\n", key, value)
	return nil
}

func maskAPIKey(key string) string {
	if key == "" {
		return "<not set>"
	}
	if len(key) < 8 {
		return "***"
	}
	return key[:4] + "****" + key[len(key)-4:]
}
