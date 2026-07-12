package main

import (
	"encoding/json"
	"testing"
)

func TestCoerceStringifiedArray(t *testing.T) {
	t.Parallel()

	coerced, ok := coerceStringifiedArray(
		json.RawMessage(`"[\"https://example.test/a\",\"https://example.test/b\"]"`),
	)
	if !ok {
		t.Fatal("coerceStringifiedArray returned ok=false")
	}

	var values []string
	if err := json.Unmarshal(coerced, &values); err != nil {
		t.Fatalf("unmarshal coerced array: %v", err)
	}
	if len(values) != 2 {
		t.Errorf("values = %v, want two URLs", values)
	}
}

func TestCoerceStringifiedArray_LeavesGenuineArrayUntouched(t *testing.T) {
	t.Parallel()

	if _, ok := coerceStringifiedArray(json.RawMessage(`["a","b"]`)); ok {
		t.Fatal("genuine array was reported as stringified")
	}
}
