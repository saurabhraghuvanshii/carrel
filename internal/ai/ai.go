// Package ai asks a language model to explain a problem or a failed run.
// The key is read from the user's local config and sent only to the provider
// they chose. With ExplainOnly set, the model is told never to write the
// full solution.
package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Provider    string // anthropic | openai | ollama
	Model       string
	APIKey      string
	ExplainOnly bool

	Statement  string
	Language   string
	Code       string
	LastReport string
	Question   string
}

const explainOnlyRule = "Never write the full solution or complete working code for the problem. " +
	"Explain concepts, point out what the learner's code does wrong, and give one hint at a time. " +
	"Short snippets that illustrate an idea are fine."

func system(explainOnly bool) string {
	s := "You are a patient tutor helping someone practise data structures and algorithms for coding interviews. "
	if explainOnly {
		s += explainOnlyRule
	} else {
		s += "You may show a full solution if the learner clearly asks for it."
	}
	return s
}

func user(r Request) string {
	var b strings.Builder
	b.WriteString("Problem:\n" + r.Statement + "\n\n")
	if strings.TrimSpace(r.Code) != "" {
		b.WriteString("My " + r.Language + " code:\n" + r.Code + "\n\n")
	}
	if r.LastReport != "" {
		b.WriteString("Result of my last run:\n" + r.LastReport + "\n\n")
	}
	q := strings.TrimSpace(r.Question)
	if q == "" {
		q = "Give me a hint about where to start."
	}
	b.WriteString("My question: " + q)
	return b.String()
}

var client = &http.Client{Timeout: 90 * time.Second}

// Ask sends the request to the chosen provider and returns the reply text.
func Ask(ctx context.Context, r Request) (string, error) {
	if strings.TrimSpace(r.Model) == "" {
		return "", errors.New("choose a model in Settings first")
	}
	switch r.Provider {
	case "anthropic":
		return askAnthropic(ctx, r)
	case "openai":
		return askOpenAI(ctx, r)
	case "ollama":
		return askOllama(ctx, r)
	}
	return "", fmt.Errorf("unknown provider %q", r.Provider)
}

func post(ctx context.Context, url string, headers map[string]string, body any, out any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("the provider answered %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	return json.Unmarshal(data, out)
}

func askAnthropic(ctx context.Context, r Request) (string, error) {
	var out struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	err := post(ctx, "https://api.anthropic.com/v1/messages",
		map[string]string{"x-api-key": r.APIKey, "anthropic-version": "2023-06-01"},
		map[string]any{
			"model":      r.Model,
			"max_tokens": 1024,
			"system":     system(r.ExplainOnly),
			"messages":   []map[string]string{{"role": "user", "content": user(r)}},
		}, &out)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	for _, c := range out.Content {
		sb.WriteString(c.Text)
	}
	return sb.String(), nil
}

func askOpenAI(ctx context.Context, r Request) (string, error) {
	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	err := post(ctx, "https://api.openai.com/v1/chat/completions",
		map[string]string{"Authorization": "Bearer " + r.APIKey},
		map[string]any{
			"model": r.Model,
			"messages": []map[string]string{
				{"role": "system", "content": system(r.ExplainOnly)},
				{"role": "user", "content": user(r)},
			},
		}, &out)
	if err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", errors.New("the provider returned no answer")
	}
	return out.Choices[0].Message.Content, nil
}

// askOllama talks to a model running on this computer. No key is needed.
func askOllama(ctx context.Context, r Request) (string, error) {
	var out struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	err := post(ctx, "http://127.0.0.1:11434/api/chat", nil,
		map[string]any{
			"model":  r.Model,
			"stream": false,
			"messages": []map[string]string{
				{"role": "system", "content": system(r.ExplainOnly)},
				{"role": "user", "content": user(r)},
			},
		}, &out)
	if err != nil {
		return "", err
	}
	return out.Message.Content, nil
}
