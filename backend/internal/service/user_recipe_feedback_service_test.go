package service

import (
	"testing"

	"menu-recommend/internal/model"
)

func TestFeedbackStatusMap(t *testing.T) {
	status := FeedbackStatusMap([]model.UserRecipeFeedback{
		{FeedbackType: "like"},
		{FeedbackType: "block"},
	})

	if !status["like"] || !status["block"] {
		t.Fatalf("status = %#v, want like and block true", status)
	}
	if status["cooked"] || status["dislike"] {
		t.Fatalf("status = %#v, want cooked and dislike false", status)
	}
}

func TestRecipeFeedbackContextScoreDelta(t *testing.T) {
	ctx := recipeFeedbackContext{
		byRecipe: map[uint]map[string]bool{
			1: {"like": true},
			2: {"dislike": true},
			3: {"cooked": true},
			4: {"block": true},
		},
	}

	if ctx.scoreDelta(1) <= 0 {
		t.Fatalf("like scoreDelta = %v, want positive", ctx.scoreDelta(1))
	}
	if ctx.scoreDelta(2) >= 0 {
		t.Fatalf("dislike scoreDelta = %v, want negative", ctx.scoreDelta(2))
	}
	if ctx.scoreDelta(3) >= 0 {
		t.Fatalf("cooked scoreDelta = %v, want negative", ctx.scoreDelta(3))
	}
	if !ctx.isBlocked(4) {
		t.Fatalf("isBlocked(4) = false, want true")
	}
}
