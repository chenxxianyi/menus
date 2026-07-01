package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"menu-recommend/internal/model"
)

type AIRecipeDraft struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	CookTime    int                  `json:"cook_time"`
	Difficulty  string               `json:"difficulty"`
	PeopleCount int                  `json:"people_count"`
	Taste       string               `json:"taste"`
	HealthTags  []string             `json:"health_tags"`
	Ingredients []AIRecipeIngredient `json:"ingredients"`
	Seasonings  []AIRecipeSeasoning  `json:"seasonings"`
	Steps       []AIRecipeStep       `json:"steps"`
	Tips        string               `json:"tips"`
	Nutrition   AIRecipeNutrition    `json:"nutrition"`
}

type AIRecipeIngredient struct {
	Name     string      `json:"name"`
	Amount   interface{} `json:"amount"`
	Unit     string      `json:"unit"`
	Emoji    string      `json:"emoji"`
	Category string      `json:"category"`
	Price    float64     `json:"price"`
}

type AIRecipeSeasoning struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
}

type AIRecipeStep struct {
	Step        int    `json:"step"`
	Description string `json:"description"`
	Tip         string `json:"tip"`
}

type AIRecipeNutrition struct {
	Calories int     `json:"calories"`
	Protein  float64 `json:"protein"`
	Fat      float64 `json:"fat"`
	Carbs    float64 `json:"carbs"`
	Fiber    float64 `json:"fiber"`
}

func (c *AIClient) GenerateRecipeDraft(ctx context.Context, dishName string) (*AIRecipeDraft, error) {
	if !c.IsConfigured() {
		return nil, ErrAIConfigMissing
	}
	dishName = strings.TrimSpace(dishName)
	if dishName == "" {
		return nil, ErrAIInvalidResponse
	}

	prompt := fmt.Sprintf(`你是家庭菜谱编写助手。请为菜品“%s”生成一份适合家庭烹饪、适合2-3人的完整菜谱。
要求：
1. 只返回 JSON，不要 Markdown，不要解释。
2. JSON 格式必须是：
{"title":"菜谱标题","description":"一句简短介绍","cook_time":45,"difficulty":"简单/中等/困难","people_count":2,"taste":"咸鲜/家常/清淡/香辣等","health_tags":["家常","高蛋白"],"ingredients":[{"name":"食材名","amount":"500","unit":"克","category":"肉蛋水产","emoji":"","price":0}],"seasonings":[{"name":"调料名","amount":"1勺"}],"steps":[{"step":1,"description":"步骤说明","tip":"可选技巧"}],"tips":"小贴士","nutrition":{"calories":420,"protein":32,"fat":18,"carbs":20,"fiber":3}}
3. ingredients 数量 4 到 15 项，seasonings 数量 2 到 10 项，steps 数量 4 到 8 步。
4. category 只能使用：肉蛋水产、蔬菜、主食、调味、配料、其他。
5. amount 和 unit 要拆开，例如 amount="500", unit="克"；如果适量，则 amount="适量", unit=""。
6. 不要编造品牌名，不要包含厨具清单。`, dishName)

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

	var draft AIRecipeDraft
	content := strings.TrimSpace(completion.Choices[0].Message.Content)
	if err := json.Unmarshal([]byte(content), &draft); err != nil {
		return nil, fmt.Errorf("%w: decode recipe content: %v", ErrAIInvalidResponse, err)
	}
	if err := normalizeAIRecipeDraft(&draft, dishName); err != nil {
		return nil, err
	}
	return &draft, nil
}

func (d *AIRecipeDraft) ToRecipe(categoryID uint) (*model.Recipe, error) {
	if d == nil {
		return nil, ErrAIInvalidResponse
	}
	healthTags, err := json.Marshal(d.HealthTags)
	if err != nil {
		return nil, err
	}
	ingredients, err := json.Marshal(d.Ingredients)
	if err != nil {
		return nil, err
	}
	seasonings, err := json.Marshal(d.Seasonings)
	if err != nil {
		return nil, err
	}
	steps, err := json.Marshal(d.Steps)
	if err != nil {
		return nil, err
	}
	nutrition, err := json.Marshal(d.Nutrition)
	if err != nil {
		return nil, err
	}

	return &model.Recipe{
		Title:       d.Title,
		Description: d.Description,
		CategoryID:  categoryID,
		CookTime:    d.CookTime,
		Difficulty:  d.Difficulty,
		PeopleCount: d.PeopleCount,
		Taste:       d.Taste,
		HealthTags:  model.JSON(healthTags),
		Ingredients: model.JSON(ingredients),
		Seasonings:  model.JSON(seasonings),
		Steps:       model.JSON(steps),
		Tips:        d.Tips,
		Nutrition:   model.JSON(nutrition),
		Status:      1,
	}, nil
}

func normalizeAIRecipeDraft(draft *AIRecipeDraft, fallbackTitle string) error {
	draft.Title = strings.TrimSpace(draft.Title)
	if draft.Title == "" {
		draft.Title = strings.TrimSpace(fallbackTitle)
	}
	draft.Description = strings.TrimSpace(draft.Description)
	if draft.Description == "" {
		draft.Description = draft.Title + "的家庭版做法，适合日常餐桌。"
	}
	if draft.CookTime <= 0 || draft.CookTime > 240 {
		draft.CookTime = 45
	}
	switch strings.TrimSpace(draft.Difficulty) {
	case "简单", "中等", "困难":
		draft.Difficulty = strings.TrimSpace(draft.Difficulty)
	default:
		draft.Difficulty = "中等"
	}
	if draft.PeopleCount <= 0 || draft.PeopleCount > 12 {
		draft.PeopleCount = 2
	}
	draft.Taste = strings.TrimSpace(draft.Taste)
	if draft.Taste == "" {
		draft.Taste = "家常"
	}

	draft.HealthTags = normalizeStringList(draft.HealthTags, 6)
	if len(draft.HealthTags) == 0 {
		draft.HealthTags = []string{"家常"}
	}

	ingredients := make([]AIRecipeIngredient, 0, len(draft.Ingredients))
	seenIngredients := map[string]bool{}
	for _, ingredient := range draft.Ingredients {
		name := strings.TrimSpace(ingredient.Name)
		if name == "" {
			continue
		}
		key := normalizeShoppingItemName(name)
		if seenIngredients[key] {
			continue
		}
		seenIngredients[key] = true
		category := normalizeShoppingCategory(ingredient.Category)
		if category == "" {
			category = inferShoppingCategory(name, "")
		}
		ingredients = append(ingredients, AIRecipeIngredient{
			Name:     name,
			Amount:   normalizeRecipeAmount(ingredient.Amount),
			Unit:     strings.TrimSpace(ingredient.Unit),
			Emoji:    strings.TrimSpace(ingredient.Emoji),
			Category: category,
			Price:    0,
		})
		if len(ingredients) >= 20 {
			break
		}
	}
	if len(ingredients) == 0 {
		return ErrAIInvalidResponse
	}
	draft.Ingredients = ingredients

	seasonings := make([]AIRecipeSeasoning, 0, len(draft.Seasonings))
	seenSeasonings := map[string]bool{}
	for _, seasoning := range draft.Seasonings {
		name := strings.TrimSpace(seasoning.Name)
		if name == "" {
			continue
		}
		key := normalizeShoppingItemName(name)
		if seenSeasonings[key] {
			continue
		}
		seenSeasonings[key] = true
		amount := strings.TrimSpace(seasoning.Amount)
		if amount == "" {
			amount = "适量"
		}
		seasonings = append(seasonings, AIRecipeSeasoning{Name: name, Amount: amount})
		if len(seasonings) >= 12 {
			break
		}
	}
	draft.Seasonings = seasonings

	steps := make([]AIRecipeStep, 0, len(draft.Steps))
	for _, step := range draft.Steps {
		description := strings.TrimSpace(step.Description)
		if description == "" {
			continue
		}
		steps = append(steps, AIRecipeStep{
			Step:        len(steps) + 1,
			Description: description,
			Tip:         strings.TrimSpace(step.Tip),
		})
		if len(steps) >= 10 {
			break
		}
	}
	if len(steps) == 0 {
		return ErrAIInvalidResponse
	}
	draft.Steps = steps

	draft.Tips = strings.TrimSpace(draft.Tips)
	if draft.Tips == "" {
		draft.Tips = "火候和调味可按家人口味微调。"
	}
	if draft.Nutrition.Calories < 0 {
		draft.Nutrition.Calories = 0
	}
	if draft.Nutrition.Protein < 0 {
		draft.Nutrition.Protein = 0
	}
	if draft.Nutrition.Fat < 0 {
		draft.Nutrition.Fat = 0
	}
	if draft.Nutrition.Carbs < 0 {
		draft.Nutrition.Carbs = 0
	}
	if draft.Nutrition.Fiber < 0 {
		draft.Nutrition.Fiber = 0
	}
	return nil
}

func normalizeStringList(values []string, limit int) []string {
	result := make([]string, 0, len(values))
	seen := map[string]bool{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
		if limit > 0 && len(result) >= limit {
			break
		}
	}
	return result
}

func normalizeRecipeAmount(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		v = strings.TrimSpace(v)
		if v == "" {
			return "适量"
		}
		return v
	case float64, int, int64:
		return v
	case nil:
		return "适量"
	default:
		text := strings.TrimSpace(fmt.Sprint(v))
		if text == "" {
			return "适量"
		}
		return text
	}
}
