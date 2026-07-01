package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type AISceneRecommendContext struct {
	Scene               string
	SceneLabel          string
	MealType            string
	PeopleCount         int
	TastePreference     []string
	HealthGoal          string
	AvoidIngredients    []string
	FavoriteIngredients []string
	CookTimePreference  string
}

type AISceneRecommendDraft struct {
	MenuName string               `json:"menu_name"`
	Reason   string               `json:"reason"`
	Dishes   []AISceneRecipeDraft `json:"dishes"`
}

type AISceneRecipeDraft struct {
	Type   string        `json:"type"`
	Reason string        `json:"reason"`
	Recipe AIRecipeDraft `json:"recipe"`
}

func (c *AIClient) GenerateSceneRecipeDrafts(ctx context.Context, scene AISceneRecommendContext) (*AISceneRecommendDraft, error) {
	if !c.IsConfigured() {
		return nil, ErrAIConfigMissing
	}

	scene.SceneLabel = strings.TrimSpace(scene.SceneLabel)
	if scene.SceneLabel == "" {
		scene.SceneLabel = "日常用餐"
	}
	scene.MealType = strings.TrimSpace(scene.MealType)
	if scene.MealType == "" {
		scene.MealType = "dinner"
	}
	if scene.PeopleCount <= 0 {
		scene.PeopleCount = 2
	}

	promptTemplate := strings.Join([]string{
		"你是家庭菜单规划和菜谱编写助手。请根据用户偏好，为“%s”场景推荐 2 到 3 道可以作为一餐搭配的新菜品，并为每道菜生成完整菜谱。",
		"用户偏好：",
		"- 餐次：%s",
		"- 人数：%d",
		"- 口味偏好：%s",
		"- 健康目标：%s",
		"- 忌口/不想吃：%s",
		"- 喜欢的食材：%s",
		"- 做饭时长偏好：%s",
		"",
		"要求：",
		"1. 可以推荐用户没吃过的新菜，不要局限于已有菜品。",
		"2. 严格避开“忌口/不想吃”里的食材，不要把它们放入食材或调料。",
		"3. 喜欢的食材可以优先使用，但不要生硬堆砌。",
		"4. 菜品要适合中国家庭厨房，食材容易购买，做法真实可操作。",
		"5. 只返回 JSON，不要 Markdown，不要解释。",
		"6. JSON 格式必须包含 menu_name、reason、dishes；每个 dish 包含 type、reason、recipe。",
		"7. recipe 必须包含 title、description、cook_time、difficulty、people_count、taste、health_tags、ingredients、seasonings、steps、tips、nutrition。",
		"8. ingredients 数量 4 到 15 项，seasonings 数量 2 到 10 项，steps 数量 4 到 8 步。",
		"9. category 只能使用：肉蛋水产、蔬菜、主食、调味、配料、其他。",
		"10. amount 和 unit 要拆开，例如 amount=\"500\", unit=\"克\"；如果适量，则 amount=\"适量\", unit=\"\"。",
		"11. 不要编造品牌名，不要包含厨具清单。",
	}, "\n")

	prompt := fmt.Sprintf(promptTemplate,
		scene.SceneLabel,
		scene.MealType,
		scene.PeopleCount,
		joinOrDefault(scene.TastePreference, "未设置"),
		emptyOrDefault(scene.HealthGoal, "未设置"),
		joinOrDefault(scene.AvoidIngredients, "无"),
		joinOrDefault(scene.FavoriteIngredients, "未设置"),
		emptyOrDefault(scene.CookTimePreference, "未设置"),
	)

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

	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
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
	var draft AISceneRecommendDraft
	if err := json.Unmarshal([]byte(content), &draft); err != nil {
		return nil, fmt.Errorf("%w: decode scene content: %v", ErrAIInvalidResponse, err)
	}
	if err := normalizeAISceneRecommendDraft(&draft, scene); err != nil {
		return nil, err
	}
	return &draft, nil
}

func normalizeAISceneRecommendDraft(draft *AISceneRecommendDraft, scene AISceneRecommendContext) error {
	draft.MenuName = strings.TrimSpace(draft.MenuName)
	if draft.MenuName == "" {
		draft.MenuName = scene.SceneLabel + " AI 菜单"
	}
	draft.Reason = strings.TrimSpace(draft.Reason)
	if draft.Reason == "" {
		draft.Reason = "AI 已根据你的场景和偏好生成新菜，并同步到菜谱库。"
	}

	dishes := make([]AISceneRecipeDraft, 0, len(draft.Dishes))
	for _, item := range draft.Dishes {
		item.Type = normalizeSceneDishType(item.Type)
		item.Reason = strings.TrimSpace(item.Reason)
		if item.Reason == "" {
			item.Reason = "符合当前场景和个人偏好。"
		}
		fallbackTitle := strings.TrimSpace(item.Recipe.Title)
		if fallbackTitle == "" {
			fallbackTitle = scene.SceneLabel + "推荐菜"
		}
		if err := normalizeAIRecipeDraft(&item.Recipe, fallbackTitle); err != nil {
			continue
		}
		dishes = append(dishes, item)
		if len(dishes) >= 3 {
			break
		}
	}
	if len(dishes) == 0 {
		return ErrAIInvalidResponse
	}
	draft.Dishes = dishes
	return nil
}

func normalizeSceneDishType(value string) string {
	switch strings.TrimSpace(value) {
	case "主菜", "配菜", "汤", "主食":
		return strings.TrimSpace(value)
	default:
		return "推荐菜"
	}
}

func joinOrDefault(values []string, fallback string) string {
	values = normalizeStringList(values, 12)
	if len(values) == 0 {
		return fallback
	}
	return strings.Join(values, "、")
}

func emptyOrDefault(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}
