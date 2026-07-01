package v1

import (
	"encoding/json"
	"testing"
)

func TestUpdateShoppingRequestDetectsEmptyItemsJSON(t *testing.T) {
	var req UpdateShoppingRequest
	if err := json.Unmarshal([]byte("{\"name\":\"我的购物清单\",\"items_json\":[]}"), &req); err != nil {
		t.Fatalf("unmarshal UpdateShoppingRequest: %v", err)
	}
	if !req.hasItemsJSON {
		t.Fatal("hasItemsJSON = false, want true")
	}
	if string(req.ItemsJSON) != "[]" {
		t.Fatalf("ItemsJSON = %s, want []", string(req.ItemsJSON))
	}
}

func TestUpdateShoppingRequestIgnoresMissingItemsJSON(t *testing.T) {
	var req UpdateShoppingRequest
	if err := json.Unmarshal([]byte("{\"name\":\"只改名称\"}"), &req); err != nil {
		t.Fatalf("unmarshal UpdateShoppingRequest: %v", err)
	}
	if req.hasItemsJSON {
		t.Fatal("hasItemsJSON = true, want false")
	}
}
