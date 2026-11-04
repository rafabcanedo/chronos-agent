package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/rafabcanedo/obsidian-golang-toolkit/api/httpclient"
)

const baseURL = "https://routellm.abacus.ai/v1"

const systemPrompt = `You are an intent parser for a personal finance system.
Today's date is: %s.

Given the user's message, classify the intent into one of the following actions:
- "pay": User wants to mark an existing bill as paid (e.g., "paguei a internet", "conta de luz paga").
- "create_expense": User incurred a one-off/ad-hoc expense (e.g., "gastei 20 no sorvete", "paguei 50 no almoço").
- "create_recurring": User wants to define a recurring bill (e.g., "cadastrar aluguel de 1200 todo dia 5").
- "query": User is asking a financial question (e.g., "quanto gastei esse mês?").

Respond ONLY with a valid JSON object matching this schema:
{
  "action": "pay" | "create_expense" | "create_recurring" | "query",
  "description": string,
  "amount": number or null,
  "date": "YYYY-MM-DD" or null,
  "billing_day": integer (1-31) or null,
  "category": string or null
}`

type Action struct {
	Action      string   `json:"action"`
	Description string   `json:"description"`
	Amount      *float64 `json:"amount"`
	Date        *string  `json:"date"`
	BillingDay  *int     `json:"billing_day"`
	Category    *string  `json:"category"`
}

type Client struct {
	http  *httpclient.Client
	model string
}

func NewClient(apiKey, model string) *Client {
	return &Client{
		http: httpclient.New(
			httpclient.WithHeader("Authorization", "Bearer "+apiKey),
			httpclient.WithHeader("Content-Type", "application/json"),
		),
		model: model,
	}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model          string        `json:"model"`
	Messages       []chatMessage `json:"messages"`
	ResponseFormat struct {
		Type string `json:"type"`
	} `json:"response_format"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (c *Client) Parse(ctx context.Context, userMessage, currentDate string) (Action, error) {
	payload := chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: fmt.Sprintf(systemPrompt, currentDate)},
			{Role: "user", Content: userMessage},
		},
	}
	payload.ResponseFormat.Type = "json_object"

	body, err := json.Marshal(payload)
	if err != nil {
		return Action{}, err
	}

	resp, err := c.http.Post(ctx, baseURL+"/chat/completions", body)
	if err != nil {
		return Action{}, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return Action{}, err
	}

	if resp.StatusCode != 200 {
		return Action{}, fmt.Errorf("abacus error %d: %s", resp.StatusCode, string(raw))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(raw, &chatResp); err != nil {
		return Action{}, err
	}

	if len(chatResp.Choices) == 0 {
		return Action{}, fmt.Errorf("empty response from LLM")
	}

	var action Action
	if err := json.Unmarshal([]byte(chatResp.Choices[0].Message.Content), &action); err != nil {
		return Action{}, err
	}

	return action, nil
}
