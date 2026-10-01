// Package ai asks a language model to explain a problem or a failed run.
// The key is read from the user's local config and sent only to the provider
// they chose. With ExplainOnly set, the model is told never to write the
// full solution, and replies with long code blocks are asked again once and
// then trimmed.
//
// Request and response shapes were checked against the providers' official
// documentation on 1 October 2026:
//   - Anthropic Messages API: POST /v1/messages, Authorization: Bearer <key>
//     (x-api-key is the documented legacy fallback), anthropic-version
//     2023-06-01; the reply is in content[].text.
//   - OpenAI Responses API (recommended over Chat Completions for new
//     projects): POST /v1/responses with instructions and input; the reply is
//     in output[] items of type "message", content[] of type "output_text".
//   - Ollama: POST /api/chat with stream false; the reply is in message.content.
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

// Endpoints are variables so tests can point them at a local fake.
var (
	anthropicURL = "https://api.anthropic.com/v1/messages"
	openaiURL    = "https://api.openai.com/v1/responses"
	ollamaURL    = "http://127.0.0.1:11434/api/chat"
)

const (
	maxCode       = 20000 // characters of the learner's code sent
	maxQuestion   = 2000
	maxReport     = 2000
	maxCodeLines  = 12 // longest code block allowed with Explain only
	replyTokens   = 1024
	codeRemoved   = "[Code removed because Explain only is on.]"
	tooMuchCodeIn = "Your last answer contained too much code. Answer again in words, with at most a few lines of code to illustrate."
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

	// Ping sends a tiny request to check the settings, ignoring the rest.
	Ping bool
}

const explainOnlyRule = "Never write the full solution or complete working code for the problem. " +
	"Explain concepts, point out what the learner's code does wrong, and give one hint at a time. " +
	"Short snippets of a few lines that illustrate an idea are fine."

func system(explainOnly bool) string {
	s := "You are a patient tutor helping someone practise data structures and algorithms for coding interviews. "
	if explainOnly {
		s += explainOnlyRule
	} else {
		s += "You may show a full solution if the learner clearly asks for it."
	}
	return s
}

func cut(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max]) + "\n[cut]"
}

// user builds the only text sent about the learner: the statement, their
// code, a short summary of the last run, and their question.
func user(r Request) string {
	var b strings.Builder
	b.WriteString("Problem:\n" + r.Statement + "\n\n")
	if strings.TrimSpace(r.Code) != "" {
		b.WriteString("My " + r.Language + " code:\n" + cut(r.Code, maxCode) + "\n\n")
	}
	if r.LastReport != "" {
		b.WriteString("Result of my last run:\n" + cut(r.LastReport, maxReport) + "\n\n")
	}
	q := strings.TrimSpace(cut(r.Question, maxQuestion))
	if q == "" {
		q = "Give me a hint about where to start."
	}
	b.WriteString("My question: " + q)
	return b.String()
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

var client = &http.Client{Timeout: 90 * time.Second}

// Ask sends the request to the chosen provider and returns the reply text.
// Errors never contain the API key.
func Ask(ctx context.Context, r Request) (string, error) {
	text, err := ask(ctx, r)
	if err != nil && r.APIKey != "" {
		err = errors.New(strings.ReplaceAll(err.Error(), r.APIKey, "[your key]"))
	}
	return text, err
}

func ask(ctx context.Context, r Request) (string, error) {
	if strings.TrimSpace(r.Model) == "" {
		return "", errors.New("choose a model in Settings first")
	}
	if r.Ping {
		text, err := chat(ctx, r, "You check that a connection works.", []message{{"user", "Reply with the single word ok."}}, 16)
		if errors.Is(err, errNoAnswer) {
			return "", nil // the provider accepted the key and model; that is all a ping checks
		}
		return text, err
	}

	msgs := []message{{"user", user(r)}}
	answer, err := chat(ctx, r, system(r.ExplainOnly), msgs, replyTokens)
	if err != nil || !r.ExplainOnly {
		return answer, err
	}
	if _, long := trimLongCode(answer); !long {
		return answer, nil
	}
	msgs = append(msgs, message{"assistant", answer}, message{"user", tooMuchCodeIn})
	again, err := chat(ctx, r, system(true), msgs, replyTokens)
	if err != nil {
		again = answer
	}
	trimmed, _ := trimLongCode(again)
	return trimmed, nil
}

// trimLongCode replaces fenced code blocks longer than maxCodeLines with a
// short note. An unclosed fence runs to the end of the text.
func trimLongCode(text string) (string, bool) {
	var out, block []string
	fence := ""
	found := false
	flush := func(closed bool) {
		lines := len(block) - 1 // without the opening fence
		if closed {
			lines-- // and without the closing one
		}
		if lines > maxCodeLines {
			out = append(out, "", codeRemoved, "")
			found = true
		} else {
			out = append(out, block...)
		}
		block = nil
	}
	for _, line := range strings.Split(text, "\n") {
		t := strings.TrimSpace(line)
		switch {
		case fence == "" && (strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~")):
			fence = t[:3]
			block = []string{line}
		case fence != "" && strings.HasPrefix(t, fence):
			block = append(block, line)
			fence = ""
			flush(true)
		case fence != "":
			block = append(block, line)
		default:
			out = append(out, line)
		}
	}
	if fence != "" {
		flush(false)
	}
	return strings.Join(out, "\n"), found
}

func chat(ctx context.Context, r Request, sys string, msgs []message, maxTokens int) (string, error) {
	switch r.Provider {
	case "anthropic":
		return askAnthropic(ctx, r, sys, msgs, maxTokens)
	case "openai":
		return askOpenAI(ctx, r, sys, msgs)
	case "ollama":
		return askOllama(ctx, r, sys, msgs)
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
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return errors.New("the provider did not answer in time")
		}
		return fmt.Errorf("could not reach the provider: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return errors.New("the provider's answer was cut off")
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("the provider answered %d: %s", resp.StatusCode, errorText(data))
	}
	if err := json.Unmarshal(data, out); err != nil {
		return errors.New("the provider sent a reply that could not be read")
	}
	return nil
}

// errorText pulls the message out of a provider's error body.
func errorText(data []byte) string {
	var e struct {
		Error json.RawMessage `json:"error"`
	}
	if json.Unmarshal(data, &e) == nil && len(e.Error) > 0 {
		var withMessage struct {
			Message string `json:"message"`
		}
		var plain string
		switch {
		case json.Unmarshal(e.Error, &withMessage) == nil && withMessage.Message != "":
			return cut(withMessage.Message, 300)
		case json.Unmarshal(e.Error, &plain) == nil && plain != "":
			return cut(plain, 300)
		}
	}
	return cut(strings.TrimSpace(string(data)), 300)
}

var errNoAnswer = errors.New("the provider returned no answer")

func askAnthropic(ctx context.Context, r Request, sys string, msgs []message, maxTokens int) (string, error) {
	var out struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	err := post(ctx, anthropicURL,
		map[string]string{"Authorization": "Bearer " + r.APIKey, "anthropic-version": "2023-06-01"},
		map[string]any{"model": r.Model, "max_tokens": maxTokens, "system": sys, "messages": msgs}, &out)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	for _, c := range out.Content {
		if c.Type == "text" {
			sb.WriteString(c.Text)
		}
	}
	if sb.Len() == 0 {
		return "", errNoAnswer
	}
	return sb.String(), nil
}

func askOpenAI(ctx context.Context, r Request, sys string, msgs []message) (string, error) {
	var out struct {
		Output []struct {
			Type    string `json:"type"`
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}
	err := post(ctx, openaiURL,
		map[string]string{"Authorization": "Bearer " + r.APIKey},
		// No max_output_tokens: reasoning models spend it on thinking first and
		// can return no text at all.
		map[string]any{"model": r.Model, "instructions": sys, "input": msgs}, &out)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	for _, item := range out.Output {
		if item.Type != "message" {
			continue
		}
		for _, c := range item.Content {
			if c.Type == "output_text" {
				sb.WriteString(c.Text)
			}
		}
	}
	if sb.Len() == 0 {
		return "", errNoAnswer
	}
	return sb.String(), nil
}

// askOllama talks to a model running on this computer. No key is needed.
func askOllama(ctx context.Context, r Request, sys string, msgs []message) (string, error) {
	var out struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	all := append([]message{{"system", sys}}, msgs...)
	err := post(ctx, ollamaURL, nil, map[string]any{"model": r.Model, "stream": false, "messages": all}, &out)
	if err != nil {
		return "", err
	}
	if out.Message.Content == "" {
		return "", errNoAnswer
	}
	return out.Message.Content, nil
}
