package service

import (
	"testing"
	"time"
)

func TestCurrentCookedStreak(t *testing.T) {
	location := time.FixedZone("CST", 8*60*60)
	today := time.Date(2026, 7, 10, 12, 0, 0, 0, location)
	stats := map[string]*DailyFoodStat{
		"2026-07-10": {CookedCount: 1},
		"2026-07-09": {CookedCount: 2},
		"2026-07-08": {CookedCount: 1},
	}
	if got := currentCookedStreak(stats, today); got != 3 {
		t.Fatalf("currentCookedStreak() = %d, want 3", got)
	}
}

func TestRecipeNutritionRejectsEmptyData(t *testing.T) {
	if _, complete := recipeNutrition([]byte(`{}`)); complete {
		t.Fatal("recipeNutrition() marked empty nutrition as complete")
	}
}
