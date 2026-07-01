package service

import (
	"errors"
	"testing"
)

func TestParseAISceneRecommendDraftAcceptsWrappedRecipes(t *testing.T) {
	content := "```json\n" + `{
		"data": {
			"menuName": "家庭聚餐 AI 菜单",
			"summary": "适合多人分享的家常搭配",
			"recipes": [
				{
					"type": "主菜",
					"reason": "适合四人以上分享",
					"dish_name": "香菇滑鸡",
					"description": "嫩滑鲜香",
					"cookTime": "35分钟",
					"difficulty": "简单",
					"servings": 4,
					"flavor": "咸鲜",
					"tags": "家常、高蛋白",
					"ingredients": ["鸡腿肉 500克", "香菇 8朵", "青椒 1个"],
					"seasonings": ["生抽 1勺", "盐 适量"],
					"steps": ["鸡腿肉切块腌制", "香菇切片", "下锅炒熟收汁"],
					"tips": "鸡肉先腌更入味"
				}
			]
		}
	}` + "\n```"

	draft, err := parseAISceneRecommendDraft(content, AISceneRecommendContext{SceneLabel: "家庭聚餐"})
	if err != nil {
		t.Fatalf("parseAISceneRecommendDraft() error = %v", err)
	}
	if draft.MenuName != "家庭聚餐 AI 菜单" {
		t.Fatalf("MenuName = %q", draft.MenuName)
	}
	if len(draft.Dishes) != 1 {
		t.Fatalf("len(Dishes) = %d, want 1", len(draft.Dishes))
	}
	recipe := draft.Dishes[0].Recipe
	if recipe.Title != "香菇滑鸡" {
		t.Fatalf("Recipe.Title = %q", recipe.Title)
	}
	if recipe.CookTime != 35 {
		t.Fatalf("CookTime = %d, want 35", recipe.CookTime)
	}
	if len(recipe.Ingredients) != 3 {
		t.Fatalf("len(Ingredients) = %d, want 3", len(recipe.Ingredients))
	}
	if len(recipe.Steps) != 3 {
		t.Fatalf("len(Steps) = %d, want 3", len(recipe.Steps))
	}
}

func TestParseAISceneRecommendDraftFallsBackWhenDishesAreFlat(t *testing.T) {
	content := `{
		"menu_name": "快手一餐",
		"dishes": [
			{
				"type": "主菜",
				"dish_name": "青椒鸡蛋",
				"cook_time": "15分钟",
				"ingredients": [{"名称":"鸡蛋","用量":"3","单位":"个"}, {"名称":"青椒","用量":"2","单位":"个"}],
				"seasonings": ["盐 适量", "生抽 半勺"],
				"steps": [{"说明":"鸡蛋打散"}, {"说明":"青椒切块"}, {"说明":"炒熟调味"}]
			}
		]
	}`

	draft, err := parseAISceneRecommendDraft(content, AISceneRecommendContext{SceneLabel: "快手一餐"})
	if err != nil {
		t.Fatalf("parseAISceneRecommendDraft() error = %v", err)
	}
	if len(draft.Dishes) != 1 {
		t.Fatalf("len(Dishes) = %d, want 1", len(draft.Dishes))
	}
	recipe := draft.Dishes[0].Recipe
	if recipe.Title != "青椒鸡蛋" {
		t.Fatalf("Recipe.Title = %q", recipe.Title)
	}
	if len(recipe.Ingredients) != 2 {
		t.Fatalf("len(Ingredients) = %d, want 2", len(recipe.Ingredients))
	}
	if len(recipe.Steps) != 3 {
		t.Fatalf("len(Steps) = %d, want 3", len(recipe.Steps))
	}
}

func TestParseAIRecipeDraftAcceptsChineseFieldNames(t *testing.T) {
	content := `这里是菜谱：{
		"菜名": "番茄牛肉汤",
		"介绍": "酸甜鲜美",
		"烹饪时间": "40分钟",
		"难度": "中等",
		"人数": 3,
		"口味": "酸甜",
		"健康标签": ["高蛋白"],
		"食材": [{"名称":"牛肉","用量":"300","单位":"克"}, {"名称":"番茄","用量":"2","单位":"个"}],
		"调料": ["盐 适量", "生抽 1勺"],
		"步骤": ["牛肉焯水", "番茄炒出汁", "加水炖煮"]
	}`

	draft, err := parseAIRecipeDraft(content, "番茄牛肉汤")
	if err != nil {
		t.Fatalf("parseAIRecipeDraft() error = %v", err)
	}
	if draft.Title != "番茄牛肉汤" {
		t.Fatalf("Title = %q", draft.Title)
	}
	if draft.CookTime != 40 {
		t.Fatalf("CookTime = %d, want 40", draft.CookTime)
	}
	if len(draft.Ingredients) != 2 {
		t.Fatalf("len(Ingredients) = %d, want 2", len(draft.Ingredients))
	}
	if draft.Ingredients[0].Name != "牛肉" {
		t.Fatalf("first ingredient = %#v", draft.Ingredients[0])
	}
}

func TestParseAIShoppingSuggestionAcceptsAlternativeKeys(t *testing.T) {
	content := `{"result":{"shopping_list":["鸡腿肉 500克","香菇 8朵","生抽 1瓶"]}}`

	suggestion, err := parseAIShoppingSuggestion(content)
	if err != nil {
		t.Fatalf("parseAIShoppingSuggestion() error = %v", err)
	}
	if len(suggestion.Items) != 3 {
		t.Fatalf("len(Items) = %d, want 3", len(suggestion.Items))
	}
	if suggestion.Items[0].Name != "鸡腿肉" {
		t.Fatalf("first item = %#v", suggestion.Items[0])
	}
}

func TestDecodeAIJSONContentRejectsPlainText(t *testing.T) {
	_, err := decodeAIJSONContent("今天推荐香菇滑鸡，做法很简单。")
	if !errors.Is(err, ErrAIInvalidResponse) {
		t.Fatalf("error = %v, want ErrAIInvalidResponse", err)
	}
}
