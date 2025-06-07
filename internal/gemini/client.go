package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

	"gdmc/internal/config"
)

type Client struct {
	apiKey string
	config *config.Config
}

type TemplateData struct {
	Commits  string
	Interval string
	Author   string
}

func NewClient(apiKey string, cfg *config.Config) *Client {
	return &Client{
		apiKey: apiKey,
		config: cfg,
	}
}

func (c *Client) GenerateMessage(commits []string, templateKey, intervalKey, author string) (string, error) {
	tmpl, exists := c.config.Templates[templateKey]
	if !exists {
		return "", fmt.Errorf("unknown template: %s", templateKey)
	}

	interval, exists := c.config.Intervals[intervalKey]
	if !exists {
		return "", fmt.Errorf("unknown interval: %s", intervalKey)
	}

	data := TemplateData{
		Commits:  strings.Join(commits, "\n"),
		Interval: interval.Name,
		Author:   author,
	}

	prompt, err := c.processTemplate(tmpl.Prompt, data)
	if err != nil {
		return "", fmt.Errorf("failed to process template: %w", err)
	}

	return c.callAPI(prompt)
}

func (c *Client) processTemplate(tmplStr string, data TemplateData) (string, error) {
	tmpl, err := template.New("prompt").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (c *Client) callAPI(prompt string) (string, error) {
	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{
						Text: prompt,
					},
				},
			},
		},
		GenerationConfig: GenerationConfig{
			Temperature:     0.7,
			TopK:            32,
			TopP:            1,
			MaxOutputTokens: 1000,
		},
		SafetySettings: []SafetySetting{
			{
				Category:  "HARM_CATEGORY_HARASSMENT",
				Threshold: "BLOCK_MEDIUM_AND_ABOVE",
			},
			{
				Category:  "HARM_CATEGORY_HATE_SPEECH",
				Threshold: "BLOCK_MEDIUM_AND_ABOVE",
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=%s", c.apiKey)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var response GeminiResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from API")
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}

func (c *Client) GetAvailableTemplates() map[string]config.Template {
	return c.config.Templates
}

type GeminiRequest struct {
	Contents         []Content        `json:"contents"`
	GenerationConfig GenerationConfig `json:"generationConfig"`
	SafetySettings   []SafetySetting  `json:"safetySettings"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GenerationConfig struct {
	Temperature     float64 `json:"temperature"`
	TopK            int     `json:"topK"`
	TopP            float64 `json:"topP"`
	MaxOutputTokens int     `json:"maxOutputTokens"`
}

type SafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content ContentResponse `json:"content"`
}

type ContentResponse struct {
	Parts []PartResponse `json:"parts"`
}

type PartResponse struct {
	Text string `json:"text"`
}
