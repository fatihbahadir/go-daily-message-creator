package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Author      string              `json:"author"`
	APIKey      string              `json:"api_key,omitempty"`
	DefaultType string              `json:"default_type"`
	Intervals   map[string]Interval `json:"intervals"`
	Templates   map[string]Template `json:"templates"`
	GitSettings GitSettings         `json:"git_settings"`
}

type Interval struct {
	Since string `json:"since"`
	Until string `json:"until"`
	Name  string `json:"name"`
}

type Template struct {
	Name        string `json:"name"`
	Prompt      string `json:"prompt"`
	Description string `json:"description"`
}

type GitSettings struct {
	IncludeMerges bool     `json:"include_merges"`
	Branches      []string `json:"branches"`
	ExcludePaths  []string `json:"exclude_paths"`
}

func Load() (*Config, error) {
	configPath := getConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return createDefaultConfig(configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
		config.APIKey = apiKey
	}

	return &config, nil
}

func (c *Config) Save() error {
	configPath := getConfigPath()

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "gdmc", "config.json")
}

func createDefaultConfig(path string) (*Config, error) {
	config := &Config{
		Author:      "",
		DefaultType: "report",
		Intervals: map[string]Interval{
			"daily": {
				Since: "yesterday.midnight",
				Until: "now",
				Name:  "Daily",
			},
			"weekly": {
				Since: "1.week.ago",
				Until: "now",
				Name:  "Weekly",
			},
			"monthly": {
				Since: "1.month.ago",
				Until: "now",
				Name:  "Monthly",
			},
		},
		Templates: map[string]Template{
			"report": {
				Name:        "Status Report",
				Description: "Professional status report format",
				Prompt: `Based on the following git commits from the {{.Interval}} period, create a professional status report:
				Git Commits:
				{{.Commits}}

				Create a structured report with:
				1. **Summary**: Brief overview of accomplishments
				2. **Key Changes**: Main features or improvements
				3. **Technical Details**: Important technical aspects
				4. **Impact**: How these changes benefit the project
				5. **Next Steps**: Planned future work

				Format as a professional status update.`,
			},
			"transcript": {
				Name:        "Meeting Transcript",
				Description: "Daily standup meeting format",
				Prompt: `Based on the following git commits from the {{.Interval}} period, create a standup meeting update:

				Git Commits:
				{{.Commits}}

				Format as a standup meeting entry:
				- **What I accomplished**: Summary of completed work
				- **Current focus**: What I'm working on now
				- **Next priorities**: Upcoming tasks
				- **Blockers/Notes**: Any challenges or important notes

				Keep it conversational and concise.`,
			},
			"summary": {
				Name:        "Work Summary",
				Description: "Concise work summary",
				Prompt: `Summarize the following git commits from the {{.Interval}} period:

				{{.Commits}}

				Provide a concise summary of the work done, highlighting the most important changes and their purpose.`,
			},
		},
		GitSettings: GitSettings{
			IncludeMerges: false,
			Branches:      []string{"--all"},
			ExcludePaths:  []string{},
		},
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := config.Save(); err != nil {
		return nil, fmt.Errorf("failed to save default config: %w", err)
	}

	return config, nil
}
