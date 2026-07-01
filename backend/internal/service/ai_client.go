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
var ErrAIUpstream = errors.New("ai upstream error")

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

	content, err := c.chatCompletion(ctx, reqBody, 1<<20)
	if err != nil {
		return nil, err
	}

	suggestion, err := parseAIShoppingSuggestion(content)
	if err != nil {
		return nil, err
	}

	items := normalizeAIShoppingItems(suggestion.Items)
	if len(items) == 0 {
		return nil, ErrAIInvalidResponse
	}
	return items, nil
}

func parseAIShoppingSuggestion(content string) (*AIShoppingSuggestion, error) {
	raw, err := decodeAIJSONContent(content)
	if err != nil {
		return nil, err
	}
	raw = unwrapAIObject(raw, "items", "items", "ingredients", "shopping_list", "采购清单", "食材")

	var suggestion AIShoppingSuggestion
	adapted := false
	if err := json.Unmarshal(raw, &suggestion); err != nil {
		if adaptedSuggestion, ok := adaptAIShoppingSuggestion(raw); ok {
			suggestion = adaptedSuggestion
			adapted = true
		} else {
			return nil, fmt.Errorf("%w: decode content: %v", ErrAIInvalidResponse, err)
		}
	} else if len(suggestion.Items) == 0 {
		if adaptedSuggestion, ok := adaptAIShoppingSuggestion(raw); ok {
			suggestion = adaptedSuggestion
			adapted = true
		}
	}
	if !adapted && len(suggestion.Items) > 0 && len(normalizeAIShoppingItems(suggestion.Items)) == 0 {
		if adaptedSuggestion, ok := adaptAIShoppingSuggestion(raw); ok {
			suggestion = adaptedSuggestion
		}
	}
	return &suggestion, nil
}

func adaptAIShoppingSuggestion(raw json.RawMessage) (AIShoppingSuggestion, bool) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return AIShoppingSuggestion{}, false
	}
	itemsRaw := rawValue(obj, "items", "ingredients", "shopping_list", "list", "采购清单", "食材")
	if len(itemsRaw) == 0 {
		return AIShoppingSuggestion{}, false
	}

	var rawItems []json.RawMessage
	if err := json.Unmarshal(itemsRaw, &rawItems); err != nil {
		if textItems := adaptAIShoppingItemsFromText(rawToString(itemsRaw)); len(textItems) > 0 {
			return AIShoppingSuggestion{Items: textItems}, true
		}
		return AIShoppingSuggestion{}, false
	}

	items := make([]DishShoppingItem, 0, len(rawItems))
	for _, rawItem := range rawItems {
		if item, ok := adaptAIShoppingItem(rawItem); ok {
			items = append(items, item)
		}
	}
	return AIShoppingSuggestion{Items: items}, len(items) > 0
}

func adaptAIShoppingItem(raw json.RawMessage) (DishShoppingItem, bool) {
	if isJSONObject(raw) {
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(raw, &obj); err != nil {
			return DishShoppingItem{}, false
		}
		name := rawString(obj, "name", "title", "ingredient", "食材", "名称")
		amount := firstNonEmpty(rawString(obj, "amount", "quantity", "qty", "用量", "数量"), "适量")
		return DishShoppingItem{
			Name:     name,
			Amount:   amount,
			Emoji:    rawString(obj, "emoji"),
			Category: rawString(obj, "category", "type", "分类", "类别"),
			Price:    rawFloat(obj, "price", "价格"),
			Checked:  false,
		}, name != ""
	}
	text := rawToString(raw)
	name, amount, _ := parseIngredientText(text)
	return DishShoppingItem{Name: name, Amount: firstNonEmpty(amount, "适量")}, name != ""
}

func adaptAIShoppingItemsFromText(text string) []DishShoppingItem {
	parts := splitAITextList(text)
	items := make([]DishShoppingItem, 0, len(parts))
	for _, part := range parts {
		name, amount, _ := parseIngredientText(part)
		if name != "" {
			items = append(items, DishShoppingItem{Name: name, Amount: firstNonEmpty(amount, "适量")})
		}
	}
	return items
}

func (c *AIClient) chatCompletion(ctx context.Context, reqBody chatCompletionRequest, bodyLimit int64) (string, error) {
	content, statusCode, body, err := c.doChatCompletion(ctx, reqBody, bodyLimit)
	if err == nil {
		return content, nil
	}
	if reqBody.ResponseFormat != nil && statusCode > 0 && isResponseFormatUnsupported(body) {
		reqBody.ResponseFormat = nil
		return c.chatCompletion(ctx, reqBody, bodyLimit)
	}
	return "", err
}

func (c *AIClient) doChatCompletion(ctx context.Context, reqBody chatCompletionRequest, bodyLimit int64) (string, int, []byte, error) {
	encoded, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(encoded))
	if err != nil {
		return "", 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", 0, nil, fmt.Errorf("%w: %v", ErrAIUpstream, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, bodyLimit))
	if err != nil {
		return "", resp.StatusCode, body, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", resp.StatusCode, body, fmt.Errorf("%w: status %d %s", ErrAIUpstream, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var completion chatCompletionResponse
	if err := json.Unmarshal(body, &completion); err != nil {
		return "", resp.StatusCode, body, fmt.Errorf("%w: decode completion: %v", ErrAIInvalidResponse, err)
	}
	if len(completion.Choices) == 0 {
		return "", resp.StatusCode, body, ErrAIInvalidResponse
	}
	content := strings.TrimSpace(chatContentAsString(completion.Choices[0].Message.Content))
	if content == "" {
		return "", resp.StatusCode, body, ErrAIInvalidResponse
	}
	return content, resp.StatusCode, body, nil
}

func isResponseFormatUnsupported(body []byte) bool {
	text := strings.ToLower(string(body))
	return strings.Contains(text, "response_format") &&
		(strings.Contains(text, "unsupported") ||
			strings.Contains(text, "not support") ||
			strings.Contains(text, "not supported") ||
			strings.Contains(text, "unrecognized") ||
			strings.Contains(text, "unknown"))
}

func chatContentAsString(raw json.RawMessage) string {
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return text
	}
	return string(raw)
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
		Message struct {
			Content json.RawMessage `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
