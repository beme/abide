package internal

import (
	"testing"
)

func TestUpdateKeyValuesInMap(t *testing.T) {
	m := map[string]interface{}{
		"A": map[string]interface{}{
			"B": 1,
			"C": map[string]interface{}{
				"B": 2,
			},
		},
		"B": 3,
	}

	newM := UpdateKeyValuesInMap("B", 0, m)
	b1 := newM["A"].(map[string]interface{})["B"].(int)
	if b1 != 0 {
		t.Fatalf("Expected 0, instead got %d.", b1)
	}
	b2 := newM["A"].(map[string]interface{})["C"].(map[string]interface{})["B"].(int)
	if b1 != 0 {
		t.Fatalf("Expected 0, instead got %d.", b2)
	}
	b3 := newM["B"].(int)
	if b3 != 0 {
		t.Fatalf("Expected 0, instead got %d.", b3)
	}
}
