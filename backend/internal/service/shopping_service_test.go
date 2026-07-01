package service

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"menu-recommend/internal/model"
)

func TestRemoveShoppingItems(t *testing.T) {
	items := model.JSON(`[
		{"name":"鸡蛋"},
		{"name":"牛肉"},
		{"name":"西红柿"},
		{"name":"酱油"}
	]`)

	updated, deletedCount, err := removeShoppingItems(items, []int{3, 1})
	if err != nil {
		t.Fatalf("removeShoppingItems() error = %v", err)
	}
	if deletedCount != 2 {
		t.Fatalf("deletedCount = %d, want 2", deletedCount)
	}

	var remaining []map[string]string
	if err := json.Unmarshal(updated, &remaining); err != nil {
		t.Fatalf("unmarshal updated items: %v", err)
	}
	want := []map[string]string{{"name": "鸡蛋"}, {"name": "西红柿"}}
	if !reflect.DeepEqual(remaining, want) {
		t.Fatalf("remaining = %#v, want %#v", remaining, want)
	}
}

func TestRemoveShoppingItemsRejectsInvalidIndices(t *testing.T) {
	items := model.JSON(`[ {"name":"鸡蛋"} ]`)
	tests := []struct {
		name    string
		indices []int
	}{
		{name: "empty", indices: nil},
		{name: "negative", indices: []int{-1}},
		{name: "out of range", indices: []int{1}},
		{name: "duplicate", indices: []int{0, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := removeShoppingItems(items, tt.indices)
			if !errors.Is(err, ErrInvalidShoppingItemIndices) {
				t.Fatalf("error = %v, want ErrInvalidShoppingItemIndices", err)
			}
		})
	}
}

func TestShoppingItemsFromRecipeUsesRealIngredients(t *testing.T) {
	raw := model.JSON(`[
		{"name":"鸡肉","amount":"500","unit":"克","category":"肉蛋水产"},
		{"name":"香菇","amount":"8","unit":"朵","category":"蔬菜"},
		{"name":"料酒","amount":"1","unit":"勺","category":"调味"}
	]`)

	items, err := shoppingItemsFromRecipe(raw)
	if err != nil {
		t.Fatalf("shoppingItemsFromRecipe() error = %v", err)
	}

	want := []DishShoppingItem{
		{Name: "鸡肉", Amount: "500克", Category: "肉蛋水产", Checked: false, Status: "pending"},
		{Name: "香菇", Amount: "8朵", Category: "蔬菜", Checked: false, Status: "pending"},
		{Name: "料酒", Amount: "1勺", Category: "调味", Checked: false, Status: "pending"},
	}
	if !reflect.DeepEqual(items, want) {
		t.Fatalf("items = %#v, want %#v", items, want)
	}
}

func TestShoppingItemsFromRecipeRejectsEmptyIngredients(t *testing.T) {
	_, err := shoppingItemsFromRecipe(model.JSON(`[]`))
	if !errors.Is(err, ErrRecipeIngredientsEmpty) {
		t.Fatalf("error = %v, want ErrRecipeIngredientsEmpty", err)
	}
}

func TestMergeShoppingItemsCombinesAliasAndSameUnit(t *testing.T) {
	items := []DishShoppingItem{
		{Name: "番茄", Amount: "2个", Category: "蔬菜"},
		{Name: "西红柿", Amount: "3个", Category: "蔬菜"},
		{Name: "鸡蛋", Amount: "2个", Category: "肉蛋水产"},
	}

	got := mergeShoppingItems(items)
	want := []DishShoppingItem{
		{Name: "番茄", Amount: "5个", Category: "蔬菜", Checked: false, Status: "pending"},
		{Name: "鸡蛋", Amount: "2个", Category: "肉蛋水产", Checked: false, Status: "pending"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("merged = %#v, want %#v", got, want)
	}
}

func TestMergeShoppingItemsKeepsIncompatibleAmounts(t *testing.T) {
	items := []DishShoppingItem{
		{Name: "酱油", Amount: "1勺", Category: "调味"},
		{Name: "酱油", Amount: "适量", Category: "调味"},
		{Name: "料酒", Amount: "1勺", Category: "调味"},
		{Name: "料酒", Amount: "20毫升", Category: "调味"},
	}

	got := mergeShoppingItems(items)
	want := []DishShoppingItem{
		{Name: "酱油", Amount: "1勺", Category: "调味", Checked: false, Status: "pending"},
		{Name: "料酒", Amount: "1勺、20毫升", Category: "调味", Checked: false, Status: "pending"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("merged = %#v, want %#v", got, want)
	}
}

func TestMergeShoppingItemsKeepsOwnedWhenAllSourcesOwned(t *testing.T) {
	items := []DishShoppingItem{
		{Name: "番茄", Amount: "2个", Category: "蔬菜", Status: "owned"},
		{Name: "西红柿", Amount: "3个", Category: "蔬菜", Status: "owned"},
	}

	got := mergeShoppingItems(items)
	want := []DishShoppingItem{
		{Name: "番茄", Amount: "5个", Category: "蔬菜", Checked: false, Status: "owned"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("merged = %#v, want %#v", got, want)
	}
}
