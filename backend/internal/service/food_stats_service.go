package service

import (
	"encoding/json"
	"sort"
	"time"

	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"
)

type FoodStatsService struct {
	feedbackRepo *repository.UserRecipeFeedbackRepo
	recipeRepo   *repository.RecipeRepo
}

type NutritionTotals struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Fat      float64 `json:"fat"`
	Carbs    float64 `json:"carbs"`
	Fiber    float64 `json:"fiber"`
}

type DailyFoodStat struct {
	Date        string          `json:"date"`
	CookedCount int             `json:"cooked_count"`
	Nutrition   NutritionTotals `json:"nutrition"`
}

type CookedRecipeStat struct {
	RecipeID uint   `json:"recipe_id"`
	Title    string `json:"title"`
	CookedAt string `json:"cooked_at"`
}

type FoodStats struct {
	Period                   string             `json:"period"`
	StartDate                string             `json:"start_date"`
	EndDate                  string             `json:"end_date"`
	CookedCount              int                `json:"cooked_count"`
	CookedDays               int                `json:"cooked_days"`
	CurrentStreak            int                `json:"current_streak"`
	Nutrition                NutritionTotals    `json:"nutrition"`
	NutritionRecordedRecipes int                `json:"nutrition_recorded_recipes"`
	NutritionCompleteness    int                `json:"nutrition_completeness"`
	Daily                    []DailyFoodStat    `json:"daily"`
	RecentCookedRecipes      []CookedRecipeStat `json:"recent_cooked_recipes"`
}

func NewFoodStatsService(feedbackRepo *repository.UserRecipeFeedbackRepo, recipeRepo *repository.RecipeRepo) *FoodStatsService {
	return &FoodStatsService{feedbackRepo: feedbackRepo, recipeRepo: recipeRepo}
}

func (s *FoodStatsService) Get(userID uint, period string, now time.Time) (*FoodStats, error) {
	days := 7
	if period == "month" {
		days = 30
	}
	if period != "month" {
		period = "week"
	}
	location := now.Location()
	today := dayStart(now, location)
	start := today.AddDate(0, 0, -(days - 1))

	result := &FoodStats{
		Period:    period,
		StartDate: start.Format("2006-01-02"),
		EndDate:   today.Format("2006-01-02"),
		Daily:     make([]DailyFoodStat, 0, days),
	}
	dailyByDate := make(map[string]*DailyFoodStat, days)
	for offset := 0; offset < days; offset++ {
		date := start.AddDate(0, 0, offset).Format("2006-01-02")
		item := DailyFoodStat{Date: date}
		result.Daily = append(result.Daily, item)
		dailyByDate[date] = &result.Daily[len(result.Daily)-1]
	}

	feedbacks, err := s.feedbackRepo.FindByUser(userID)
	if err != nil {
		return nil, err
	}
	cookedFeedbacks := make([]model.UserRecipeFeedback, 0)
	ids := make([]uint, 0)
	for _, feedback := range feedbacks {
		if feedback.FeedbackType != "cooked" || feedback.UpdatedAt.Before(start) {
			continue
		}
		cookedFeedbacks = append(cookedFeedbacks, feedback)
		ids = append(ids, feedback.RecipeID)
	}
	if len(cookedFeedbacks) == 0 {
		return result, nil
	}
	recipes, err := s.recipeRepo.FindByIDs(ids)
	if err != nil {
		return nil, err
	}
	recipeByID := make(map[uint]model.Recipe, len(recipes))
	for _, recipe := range recipes {
		recipeByID[recipe.ID] = recipe
	}

	for _, feedback := range cookedFeedbacks {
		date := feedback.UpdatedAt.In(location).Format("2006-01-02")
		daily := dailyByDate[date]
		if daily == nil {
			continue
		}
		result.CookedCount++
		daily.CookedCount++
		recipe, ok := recipeByID[feedback.RecipeID]
		if !ok {
			continue
		}
		if nutrition, complete := recipeNutrition(recipe.Nutrition); complete {
			result.NutritionRecordedRecipes++
			addNutrition(&result.Nutrition, nutrition)
			addNutrition(&daily.Nutrition, nutrition)
		}
		result.RecentCookedRecipes = append(result.RecentCookedRecipes, CookedRecipeStat{
			RecipeID: recipe.ID,
			Title:    recipe.Title,
			CookedAt: feedback.UpdatedAt.In(location).Format("2006-01-02 15:04"),
		})
	}

	for _, item := range result.Daily {
		if item.CookedCount > 0 {
			result.CookedDays++
		}
	}
	result.CurrentStreak = currentCookedStreak(dailyByDate, today)
	if result.CookedCount > 0 {
		result.NutritionCompleteness = result.NutritionRecordedRecipes * 100 / result.CookedCount
	}
	sort.Slice(result.RecentCookedRecipes, func(i, j int) bool {
		return result.RecentCookedRecipes[i].CookedAt > result.RecentCookedRecipes[j].CookedAt
	})
	if len(result.RecentCookedRecipes) > 5 {
		result.RecentCookedRecipes = result.RecentCookedRecipes[:5]
	}
	return result, nil
}

func dayStart(value time.Time, location *time.Location) time.Time {
	local := value.In(location)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, location)
}

func currentCookedStreak(dailyByDate map[string]*DailyFoodStat, today time.Time) int {
	streak := 0
	for cursor := today; ; cursor = cursor.AddDate(0, 0, -1) {
		item := dailyByDate[cursor.Format("2006-01-02")]
		if item == nil || item.CookedCount == 0 {
			break
		}
		streak++
	}
	return streak
}

func recipeNutrition(raw model.JSON) (NutritionTotals, bool) {
	var values map[string]float64
	if err := json.Unmarshal(raw, &values); err != nil {
		return NutritionTotals{}, false
	}
	totals := NutritionTotals{
		Calories: values["calories"],
		Protein:  values["protein"],
		Fat:      values["fat"],
		Carbs:    values["carbs"],
		Fiber:    values["fiber"],
	}
	return totals, totals.Calories > 0 || totals.Protein > 0 || totals.Fat > 0 || totals.Carbs > 0 || totals.Fiber > 0
}

func addNutrition(target *NutritionTotals, source NutritionTotals) {
	target.Calories += source.Calories
	target.Protein += source.Protein
	target.Fat += source.Fat
	target.Carbs += source.Carbs
	target.Fiber += source.Fiber
}
