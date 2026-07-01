package service

import (
	"context"
	"encoding/json"
	"fmt"
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

	content, err := c.chatCompletion(ctx, reqBody, 1<<20)
	if err != nil {
		return nil, err
	}
	draft, err := parseAIRecipeDraft(content, dishName)
	if err != nil {
		return nil, err
	}
	return draft, nil
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

func parseAIRecipeDraft(content string, fallbackTitle string) (*AIRecipeDraft, error) {
	raw, err := decodeAIJSONContent(content)
	if err != nil {
		return nil, err
	}
	raw = unwrapAIObject(raw, "recipes", "title", "name", "ingredients", "steps", "recipe")

	var draft AIRecipeDraft
	adapted := false
	if err := json.Unmarshal(raw, &draft); err != nil {
		if adaptedDraft, ok := adaptAIRecipeDraft(raw, fallbackTitle); ok {
			draft = adaptedDraft
			adapted = true
		} else {
			return nil, fmt.Errorf("%w: decode recipe content: %v", ErrAIInvalidResponse, err)
		}
	} else if draft.Title == "" || len(draft.Ingredients) == 0 || len(draft.Steps) == 0 {
		if adaptedDraft, ok := adaptAIRecipeDraft(raw, fallbackTitle); ok {
			draft = adaptedDraft
			adapted = true
		}
	}
	if err := normalizeAIRecipeDraft(&draft, fallbackTitle); err != nil {
		if !adapted {
			if adaptedDraft, ok := adaptAIRecipeDraft(raw, fallbackTitle); ok {
				if normalizeErr := normalizeAIRecipeDraft(&adaptedDraft, fallbackTitle); normalizeErr == nil {
					return &adaptedDraft, nil
				}
			}
		}
		return nil, err
	}
	return &draft, nil
}

func adaptAIRecipeDraft(raw json.RawMessage, fallbackTitle string) (AIRecipeDraft, bool) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return AIRecipeDraft{}, false
	}

	if recipeRaw := rawValue(obj, "recipe", "recipe_detail", "detail", "菜谱", "做法"); len(recipeRaw) > 0 && isJSONObject(recipeRaw) {
		var nested map[string]json.RawMessage
		if err := json.Unmarshal(recipeRaw, &nested); err == nil && hasAnyRawKey(nested, "title", "name", "ingredients", "steps") {
			obj = nested
			raw = recipeRaw
		}
	}

	draft := AIRecipeDraft{
		Title:       firstNonEmpty(rawString(obj, "title", "name", "dish_name", "recipe_name", "菜名", "菜谱名称", "菜品名称"), fallbackTitle),
		Description: rawString(obj, "description", "desc", "intro", "summary", "介绍", "描述"),
		CookTime:    rawInt(obj, "cook_time", "cookTime", "cooking_time", "time", "烹饪时间", "制作时间"),
		Difficulty:  rawString(obj, "difficulty", "level", "难度"),
		PeopleCount: rawInt(obj, "people_count", "peopleCount", "servings", "portion", "人数", "份量"),
		Taste:       rawString(obj, "taste", "flavor", "口味"),
		HealthTags:  rawStringList(obj, 6, "health_tags", "healthTags", "tags", "nutrition_tags", "健康标签", "标签"),
		Tips:        rawString(obj, "tips", "tip", "notes", "小贴士", "提示"),
		Nutrition: AIRecipeNutrition{
			Calories: rawIntFromNested(obj, []string{"nutrition", "营养"}, "calories", "kcal", "热量"),
			Protein:  rawFloatFromNested(obj, []string{"nutrition", "营养"}, "protein", "蛋白质"),
			Fat:      rawFloatFromNested(obj, []string{"nutrition", "营养"}, "fat", "脂肪"),
			Carbs:    rawFloatFromNested(obj, []string{"nutrition", "营养"}, "carbs", "carbohydrate", "碳水", "碳水化合物"),
			Fiber:    rawFloatFromNested(obj, []string{"nutrition", "营养"}, "fiber", "膳食纤维"),
		},
	}

	draft.Ingredients = adaptAIRecipeIngredients(rawValue(obj, "ingredients", "ingredient", "食材", "主料", "食材清单"))
	draft.Seasonings = adaptAIRecipeSeasonings(rawValue(obj, "seasonings", "seasoning", "condiments", "调料", "调味料", "辅料"))
	draft.Steps = adaptAIRecipeSteps(rawValue(obj, "steps", "instructions", "method", "directions", "做法", "步骤"))

	return draft, draft.Title != "" || len(draft.Ingredients) > 0 || len(draft.Steps) > 0 || len(raw) > 0
}

func rawIntFromNested(obj map[string]json.RawMessage, parents []string, keys ...string) int {
	for _, parent := range parents {
		if nested := rawValue(obj, parent); len(nested) > 0 {
			var nestedObj map[string]json.RawMessage
			if err := json.Unmarshal(nested, &nestedObj); err == nil {
				if value := rawInt(nestedObj, keys...); value != 0 {
					return value
				}
			}
		}
	}
	return rawInt(obj, keys...)
}

func rawFloatFromNested(obj map[string]json.RawMessage, parents []string, keys ...string) float64 {
	for _, parent := range parents {
		if nested := rawValue(obj, parent); len(nested) > 0 {
			var nestedObj map[string]json.RawMessage
			if err := json.Unmarshal(nested, &nestedObj); err == nil {
				if value := rawFloat(nestedObj, keys...); value != 0 {
					return value
				}
			}
		}
	}
	return rawFloat(obj, keys...)
}

func adaptAIRecipeIngredients(raw json.RawMessage) []AIRecipeIngredient {
	raw = bytesTrimSpace(raw)
	if len(raw) == 0 {
		return nil
	}
	if isJSONArray(raw) {
		var items []json.RawMessage
		if err := json.Unmarshal(raw, &items); err != nil {
			return nil
		}
		result := make([]AIRecipeIngredient, 0, len(items))
		for _, item := range items {
			if ingredient, ok := adaptAIRecipeIngredient(item); ok {
				result = append(result, ingredient)
			}
		}
		return result
	}
	return adaptAIRecipeIngredientsFromText(rawToString(raw))
}

func adaptAIRecipeIngredient(raw json.RawMessage) (AIRecipeIngredient, bool) {
	if isJSONObject(raw) {
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(raw, &obj); err != nil {
			return AIRecipeIngredient{}, false
		}
		name := rawString(obj, "name", "title", "ingredient", "食材", "名称")
		amount := firstNonEmpty(rawString(obj, "amount", "quantity", "qty", "用量", "数量"), "适量")
		unit := rawString(obj, "unit", "单位")
		return AIRecipeIngredient{
			Name:     name,
			Amount:   amount,
			Unit:     unit,
			Emoji:    rawString(obj, "emoji"),
			Category: rawString(obj, "category", "type", "分类", "类别"),
			Price:    rawFloat(obj, "price", "价格"),
		}, name != ""
	}
	text := rawToString(raw)
	name, amount, unit := parseIngredientText(text)
	return AIRecipeIngredient{Name: name, Amount: amount, Unit: unit}, name != ""
}

func adaptAIRecipeIngredientsFromText(text string) []AIRecipeIngredient {
	parts := splitAITextList(text)
	result := make([]AIRecipeIngredient, 0, len(parts))
	for _, part := range parts {
		name, amount, unit := parseIngredientText(part)
		if name != "" {
			result = append(result, AIRecipeIngredient{Name: name, Amount: amount, Unit: unit})
		}
	}
	return result
}

func parseIngredientText(text string) (string, string, string) {
	text = cleanAIListItem(text)
	if text == "" {
		return "", "", ""
	}
	text = strings.ReplaceAll(text, "：", ":")
	text = strings.ReplaceAll(text, "，", " ")
	text = strings.ReplaceAll(text, "、", " ")
	parts := strings.SplitN(text, ":", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), ""
	}
	fields := strings.Fields(text)
	if len(fields) >= 2 {
		return fields[0], strings.Join(fields[1:], ""), ""
	}
	return text, "适量", ""
}

func adaptAIRecipeSeasonings(raw json.RawMessage) []AIRecipeSeasoning {
	raw = bytesTrimSpace(raw)
	if len(raw) == 0 {
		return nil
	}
	if isJSONArray(raw) {
		var items []json.RawMessage
		if err := json.Unmarshal(raw, &items); err != nil {
			return nil
		}
		result := make([]AIRecipeSeasoning, 0, len(items))
		for _, item := range items {
			if seasoning, ok := adaptAIRecipeSeasoning(item); ok {
				result = append(result, seasoning)
			}
		}
		return result
	}
	parts := splitAITextList(rawToString(raw))
	result := make([]AIRecipeSeasoning, 0, len(parts))
	for _, part := range parts {
		name, amount, _ := parseIngredientText(part)
		if name != "" {
			result = append(result, AIRecipeSeasoning{Name: name, Amount: firstNonEmpty(amount, "适量")})
		}
	}
	return result
}

func adaptAIRecipeSeasoning(raw json.RawMessage) (AIRecipeSeasoning, bool) {
	if isJSONObject(raw) {
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(raw, &obj); err != nil {
			return AIRecipeSeasoning{}, false
		}
		name := rawString(obj, "name", "title", "seasoning", "调料", "名称")
		amount := firstNonEmpty(rawString(obj, "amount", "quantity", "qty", "用量", "数量"), "适量")
		return AIRecipeSeasoning{Name: name, Amount: amount}, name != ""
	}
	text := rawToString(raw)
	name, amount, _ := parseIngredientText(text)
	return AIRecipeSeasoning{Name: name, Amount: firstNonEmpty(amount, "适量")}, name != ""
}

func adaptAIRecipeSteps(raw json.RawMessage) []AIRecipeStep {
	raw = bytesTrimSpace(raw)
	if len(raw) == 0 {
		return nil
	}
	if isJSONArray(raw) {
		var items []json.RawMessage
		if err := json.Unmarshal(raw, &items); err != nil {
			return nil
		}
		result := make([]AIRecipeStep, 0, len(items))
		for _, item := range items {
			if step, ok := adaptAIRecipeStep(item, len(result)+1); ok {
				result = append(result, step)
			}
		}
		return result
	}
	parts := splitAITextList(rawToString(raw))
	result := make([]AIRecipeStep, 0, len(parts))
	for _, part := range parts {
		part = cleanAIListItem(part)
		if part != "" {
			result = append(result, AIRecipeStep{Step: len(result) + 1, Description: part})
		}
	}
	return result
}

func adaptAIRecipeStep(raw json.RawMessage, fallbackStep int) (AIRecipeStep, bool) {
	if isJSONObject(raw) {
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(raw, &obj); err != nil {
			return AIRecipeStep{}, false
		}
		description := rawString(obj, "description", "text", "content", "instruction", "做法", "步骤", "说明")
		step := rawInt(obj, "step", "index", "order", "序号")
		if step <= 0 {
			step = fallbackStep
		}
		return AIRecipeStep{
			Step:        step,
			Description: description,
			Tip:         rawString(obj, "tip", "tips", "提示", "小贴士"),
		}, description != ""
	}
	description := cleanAIListItem(rawToString(raw))
	return AIRecipeStep{Step: fallbackStep, Description: description}, description != ""
}

func bytesTrimSpace(raw json.RawMessage) json.RawMessage {
	return json.RawMessage(strings.TrimSpace(string(raw)))
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
