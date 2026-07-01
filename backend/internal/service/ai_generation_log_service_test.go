package service

import (
	"testing"
	"time"
)

func TestMarshalLogJSON(t *testing.T) {
	raw := marshalLogJSON(map[string]interface{}{"dish_name": "番茄炒蛋", "api_key": ""})
	if len(raw) == 0 {
		t.Fatal("marshalLogJSON() returned empty JSON")
	}
	if string(raw) != `{"api_key":"","dish_name":"番茄炒蛋"}` && string(raw) != `{"dish_name":"番茄炒蛋","api_key":""}` {
		t.Fatalf("marshalLogJSON() = %s", raw)
	}
}

func TestAIGenerationLogPayloadDurationMilliseconds(t *testing.T) {
	payload := AIGenerationLogPayload{Duration: 1500 * time.Millisecond}
	if payload.Duration.Milliseconds() != 1500 {
		t.Fatalf("Duration.Milliseconds() = %d", payload.Duration.Milliseconds())
	}
}
