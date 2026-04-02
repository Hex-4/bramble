package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func newWebSearch(workspaceDir string) Tool {
	return Tool{
		Name:        "web_search",
		Description: "Search the web for information",
		Emoji:       "🌐",
		Parameters: map[string]Parameter{
			"query": {
				Type:        "string",
				Description: "The search query",
				Required:    true,
			},
			"count": {
				Type:        "integer",
				Description: "The number of results to return. Default 5",
				Required:    false,
			},
		},
		DetailParam: "query",
		Execute: func(args map[string]any) (string, error) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			numResults := 5 // default
			if n, ok := args["count"].(float64); ok {
				numResults = int(n)
			}
			body := map[string]any{
				"query":      ArgString(args, "query"),
				"numResults": numResults,
				"type":       "auto",
				"contents": map[string]any{
					"highlights": map[string]any{
						"maxCharacters": 2000,
					},
				},
			}
			jsonBody, err := json.Marshal(body)
			if err != nil {
				return "Tool call error: failed to build request: " + err.Error(), nil
			}

			req, err := http.NewRequestWithContext(ctx, "POST", "https://api.exa.ai/search", bytes.NewReader(jsonBody))
			if err != nil {
				return "Tool call error: failed to create request: " + err.Error(), nil
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("x-api-key", os.Getenv("EXA_KEY"))

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return "Tool call error: failed to send request: " + err.Error(), nil
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Sprintf("Tool call error: exa returned status %d", resp.StatusCode), nil
			}

			var result struct {
				Results []struct {
					Title      string   `json:"title"`
					URL        string   `json:"url"`
					Highlights []string `json:"highlights"`
				} `json:"results"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return "Tool call error: failed to decode response: " + err.Error(), nil
			}

			var sb strings.Builder
			for _, r := range result.Results {
				sb.WriteString(fmt.Sprintf("## %s\n(%s)\n\n", r.Title, r.URL))
				for _, h := range r.Highlights {
					sb.WriteString(fmt.Sprintf("- %s\n\n", h))
				}
			}
			return sb.String(), nil
		},
	}
}
