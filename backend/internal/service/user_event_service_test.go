package service

import (
	"encoding/json"
	"testing"

	"menu-recommend/internal/model"
)

func TestUserEventPayloadJSONFallback(t *testing.T) {
	payload, err := json.Marshal(map[string]interface{}{"source": "detail"})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	event := model.UserEvent{PayloadJSON: model.JSON(payload)}
	if string(event.PayloadJSON) != `{"source":"detail"}` {
		t.Fatalf("PayloadJSON = %s", string(event.PayloadJSON))
	}
}
