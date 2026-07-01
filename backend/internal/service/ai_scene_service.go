package service

import (
	"context"
	"encoding/json"
	"fmt"
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

	content, err := c.chatCompletion(ctx, reqBody, 2<<20)
	if err != nil {
		return nil, err
	}
	draft, err := parseAISceneRecommendDraft(content, scene)
	if err != nil {
		return nil, err
	}
	return draft, nil
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

func parseAISceneRecommendDraft(content string, scene AISceneRecommendContext) (*AISceneRecommendDraft, error) {
	raw, err := decodeAIJSONContent(content)
	if err != nil {
		return nil, err
	}
	raw = unwrapAIObject(raw, "dishes", "dishes", "recipes", "items", "menu_name", "menuName", "菜品")

	var draft AISceneRecommendDraft
	adapted := false
	if err := json.Unmarshal(raw, &draft); err != nil {
		if adaptedDraft, ok := adaptAISceneRecommendDraft(raw, scene); ok {
			draft = adaptedDraft
			adapted = true
		} else {
			return nil, fmt.Errorf("%w: decode scene content: %v", ErrAIInvalidResponse, err)
		}
	} else if len(draft.Dishes) == 0 {
		if adaptedDraft, ok := adaptAISceneRecommendDraft(raw, scene); ok {
			draft = adaptedDraft
			adapted = true
		}
	}
	if err := normalizeAISceneRecommendDraft(&draft, scene); err != nil {
		if !adapted {
			if adaptedDraft, ok := adaptAISceneRecommendDraft(raw, scene); ok {
				if normalizeErr := normalizeAISceneRecommendDraft(&adaptedDraft, scene); normalizeErr == nil {
					return &adaptedDraft, nil
				}
			}
		}
		return nil, err
	}
	return &draft, nil
}

func adaptAISceneRecommendDraft(raw json.RawMessage, scene AISceneRecommendContext) (AISceneRecommendDraft, bool) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return AISceneRecommendDraft{}, false
	}

	draft := AISceneRecommendDraft{
		MenuName: firstNonEmpty(
			rawString(obj, "menu_name", "menuName", "name", "title", "菜单名称", "推荐菜单"),
			scene.SceneLabel+" AI 菜单",
		),
		Reason: firstNonEmpty(rawString(obj, "reason", "summary", "description", "推荐理由", "说明")),
	}

	dishRaw := rawValue(obj, "dishes", "recipes", "items", "menu", "菜品", "菜谱", "推荐菜品")
	if len(dishRaw) == 0 {
		return draft, false
	}
	if isJSONObject(dishRaw) {
		var dishObj map[string]json.RawMessage
		if err := json.Unmarshal(dishRaw, &dishObj); err == nil {
			for _, key := range []string{"dishes", "recipes", "items", "菜品", "菜谱"} {
				if nested := rawValue(dishObj, key); len(nested) > 0 {
					dishRaw = nested
					break
				}
			}
		}
	}

	var rawItems []json.RawMessage
	if err := json.Unmarshal(dishRaw, &rawItems); err != nil {
		return draft, false
	}
	for _, rawItem := range rawItems {
		if dish, ok := adaptAISceneRecipeDraft(rawItem, scene); ok {
			draft.Dishes = append(draft.Dishes, dish)
		}
	}
	return draft, len(draft.Dishes) > 0
}

func adaptAISceneRecipeDraft(raw json.RawMessage, scene AISceneRecommendContext) (AISceneRecipeDraft, bool) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return AISceneRecipeDraft{}, false
	}
	dish := AISceneRecipeDraft{
		Type:   rawString(obj, "type", "dish_type", "category", "类型", "菜品类型"),
		Reason: rawString(obj, "reason", "why", "description", "推荐理由", "说明"),
	}

	recipeRaw := rawValue(obj, "recipe", "recipe_detail", "detail", "菜谱", "做法")
	if len(recipeRaw) == 0 {
		recipeRaw = raw
	}
	recipe, ok := adaptAIRecipeDraft(recipeRaw, fallbackAISceneRecipeTitle(obj, scene))
	if !ok {
		return AISceneRecipeDraft{}, false
	}
	dish.Recipe = recipe
	return dish, true
}

func fallbackAISceneRecipeTitle(obj map[string]json.RawMessage, scene AISceneRecommendContext) string {
	return firstNonEmpty(
		rawString(obj, "title", "name", "dish_name", "recipe_name", "菜名", "菜品名称", "菜谱名称"),
		scene.SceneLabel+"推荐菜",
	)
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
