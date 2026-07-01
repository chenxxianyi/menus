package service

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

	"menu-recommend/config"
)

var ErrAIConfigMissing = errors.New("ai config missing")
var ErrAIInvalidResponse = errors.New("ai invalid response")

type AIClient struct {
	baseURL     string
	apiKey      string
	model       string
	temperature float64
	httpClient  *http.Client
}

type AIShoppingSuggestion struct {
	Items []DishShoppingItem `json:"items"`
}

func NewAIClient(cfg config.AIConfig) *AIClient {
	timeout := time.Duration(cfg.TimeoutSecs) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	temperature := cfg.Temperature
	if temperature == 0 {
		temperature = 0.2
	}
	return &AIClient{
		baseURL:     strings.TrimRight(cfg.BaseURL, "/"),
		apiKey:      strings.TrimSpace(cfg.APIKey),
		model:       strings.TrimSpace(cfg.Model),
		temperature: temperature,
		httpClient:  &http.Client{Timeout: timeout},
	}
}

func (c *AIClient) IsConfigured() bool {
	return c != nil && c.baseURL != "" && c.apiKey != "" && c.model != ""
}

func (c *AIClient) SuggestShoppingItems(ctx context.Context, dishName string) ([]DishShoppingItem, error) {
	if !c.IsConfigured() {
		return nil, ErrAIConfigMissing
	}
	dishName = strings.TrimSpace(dishName)
	if dishName == "" {
		return nil, ErrAIInvalidResponse
	}

	prompt := fmt.Sprintf(`你是家庭做饭采购助手。请根据菜品“%s”生成一份适合2-3人家庭烹饪的采购清单。
要求：
1. 只返回 JSON，不要 Markdown，不要解释。
2. JSON 格式必须是 {"items":[{"name":"食材名","amount":"数量和单位","category":"分类","emoji":"","price":0,"checked":false}]}。
3. items 数量控制在 4 到 15 项。
4. category 只能使用：肉蛋水产、蔬菜、主食、调味、配料、其他。
5. amount 必须是可读采购量，例如“500克”“8朵”“1块”“适量”。
6. 不要包含厨具、做法步骤、营养说明。`, dishName)

	reqBody := chatCompletionRequest{
		Model:       c.model,
		Temperature: c.temperature,
		Messages: []chatMessage{
			{Role: "system", Content: "你只输出严格 JSON。"},
			{Role: "user", Content: prompt},
		},
		ResponseFormat: &responseFormat{Type: "json_object"},
	}

	encoded, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(encoded))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%w: status %d %s", ErrAIInvalidResponse, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var completion chatCompletionResponse
	if err := json.Unmarshal(body, &completion); err != nil {
		return nil, fmt.Errorf("%w: decode completion: %v", ErrAIInvalidResponse, err)
	}
	if len(completion.Choices) == 0 {
		return nil, ErrAIInvalidResponse
	}

	content := strings.TrimSpace(completion.Choices[0].Message.Content)
	var suggestion AIShoppingSuggestion
	if err := json.Unmarshal([]byte(content), &suggestion); err != nil {
		return nil, fmt.Errorf("%w: decode content: %v", ErrAIInvalidResponse, err)
	}

	items := normalizeAIShoppingItems(suggestion.Items)
	if len(items) == 0 {
		return nil, ErrAIInvalidResponse
	}
	return items, nil
}

func normalizeAIShoppingItems(items []DishShoppingItem) []DishShoppingItem {
	result := make([]DishShoppingItem, 0, len(items))
	seen := make(map[string]bool)
	for _, item := range items {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}
		key := normalizeShoppingItemName(name)
		if seen[key] {
			continue
		}
		seen[key] = true
		amount := strings.TrimSpace(item.Amount)
		if amount == "" {
			amount = "适量"
		}
		category := normalizeShoppingCategory(item.Category)
		if category == "" {
			category = inferShoppingCategory(name, "")
		}
		result = append(result, DishShoppingItem{
			Name:     name,
			Amount:   amount,
			Emoji:    strings.TrimSpace(item.Emoji),
			Category: category,
			Price:    0,
			Checked:  false,
		})
		if len(result) >= 20 {
			break
		}
	}
	return result
}

func normalizeShoppingCategory(category string) string {
	category = strings.TrimSpace(category)
	switch category {
	case "肉蛋水产", "蔬菜", "主食", "调味", "配料", "其他":
		return category
	default:
		return ""
	}
}

type chatCompletionRequest struct {
	Model          string          `json:"model"`
	Messages       []chatMessage   `json:"messages"`
	Temperature    float64         `json:"temperature,omitempty"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}
