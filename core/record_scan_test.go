package core

import (
	"testing"

	"github.com/pocketbase/pocketbase/tools/types"
)

func TestNewRecordFromScannedValueMapWithJSONBytes(t *testing.T) {
	t.Parallel()

	collection := NewBaseCollection("demo_json_scan")
	collection.Fields = NewFieldsList(
		&TextField{Name: "id", PrimaryKey: true},
		&JSONField{Name: "json_object"},
	)

	record, err := newRecordFromScannedValueMap(collection, map[string]any{
		"id":          "i9naidtvr6qsgb4",
		"json_object": []byte(`{"a":{"b":"test"}}`),
	})
	if err != nil {
		t.Fatal(err)
	}

	raw, ok := record.GetRaw("json_object").(types.JSONRaw)
	if !ok {
		t.Fatalf("expected json_object to be types.JSONRaw, got %T", record.GetRaw("json_object"))
	}

	if raw.String() != `{"a":{"b":"test"}}` {
		t.Fatalf("expected json_object to round-trip, got %s", raw.String())
	}
}
