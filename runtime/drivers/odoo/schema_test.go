package odoo

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/stretchr/testify/require"
)

func TestColumnsForMapsTypesAndSplitsMany2one(t *testing.T) {
	info := map[string]fieldInfo{
		"id": {"integer"}, "name": {"char"}, "amount": {"monetary"}, "active": {"boolean"},
		"day": {"date"}, "at": {"datetime"}, "partner_id": {"many2one"}, "state": {"selection"},
	}
	cols, err := columnsFor([]string{"id", "name", "amount", "active", "day", "at", "partner_id", "state"}, info)
	require.NoError(t, err)

	s := schemaFor(cols)
	var names []string
	for _, f := range s.Fields() {
		names = append(names, f.Name)
	}
	require.Equal(t, []string{"id", "name", "amount", "active", "day", "at", "partner_id", "partner_id_name", "state"}, names)
	want := []arrow.DataType{
		arrow.PrimitiveTypes.Int64, arrow.BinaryTypes.String, arrow.PrimitiveTypes.Float64,
		arrow.FixedWidthTypes.Boolean, arrow.FixedWidthTypes.Date32, utcMicros,
		arrow.PrimitiveTypes.Int64, arrow.BinaryTypes.String, arrow.BinaryTypes.String,
	}
	for i, typ := range want {
		require.True(t, arrow.TypeEqual(typ, s.Field(i).Type), "column %s", s.Field(i).Name)
	}
}

func TestColumnsForRejectsListsAndUnknownFields(t *testing.T) {
	_, err := columnsFor([]string{"line_ids"}, map[string]fieldInfo{"line_ids": {"one2many"}})
	require.ErrorContains(t, err, `field "line_ids" has type "one2many"`)
	_, err = columnsFor([]string{"nope"}, map[string]fieldInfo{})
	require.ErrorContains(t, err, `no field "nope", or the API user cannot read it`)
}

// A field genuinely missing from fields_get and a field the API user cannot read look identical
// from the driver's side (fields_get silently omits fields gated by record rules or groups), so
// the error must not claim the field does not exist.
func TestColumnsForRejectsCollisionWithMany2oneNameColumn(t *testing.T) {
	info := map[string]fieldInfo{
		"partner_id":      {"many2one"},
		"partner_id_name": {"char"},
	}
	_, err := columnsFor([]string{"partner_id", "partner_id_name"}, info)
	require.ErrorContains(t, err, `partner_id_name`)

	// Order must not matter: the same collision, requested the other way round.
	_, err = columnsFor([]string{"partner_id_name", "partner_id"}, info)
	require.ErrorContains(t, err, `partner_id_name`)
}

func TestAppendValue(t *testing.T) {
	cases := []struct {
		kind string
		part int
		in   any
		want any // nil means the cell must be null
	}{
		{"integer", 0, json.Number("7"), int64(7)},
		{"integer", 0, false, nil},
		{"monetary", 0, json.Number("1234.5"), 1234.5},
		{"char", 0, "Acme", "Acme"},
		{"char", 0, false, nil},
		{"selection", 0, "posted", "posted"},
		{"boolean", 0, false, false},
		{"boolean", 0, true, true},
		{"date", 0, "2026-08-31", "2026-08-31"},
		{"date", 0, false, nil},
		{"datetime", 0, "2026-08-31 08:15:00", "2026-08-31 08:15:00"},
		{"datetime", 0, false, nil},
		{"many2one", 0, []any{json.Number("12"), "Cliente SL"}, int64(12)},
		{"many2one", 1, []any{json.Number("12"), "Cliente SL"}, "Cliente SL"},
		{"many2one", 1, false, nil},
		// A related record whose _rec_name is empty comes back as [id, false]: the id is still
		// readable, the display name is null instead of aborting the whole load.
		{"many2one", 0, []any{json.Number("12"), false}, int64(12)},
		{"many2one", 1, []any{json.Number("12"), false}, nil},
	}
	for _, tc := range cases {
		cols, err := columnsFor([]string{"f"}, map[string]fieldInfo{"f": {tc.kind}})
		require.NoError(t, err)
		col := cols[tc.part]
		b := array.NewBuilder(memory.DefaultAllocator, col.typ)
		require.NoError(t, col.appendValue(b, tc.in), "%s %v", tc.kind, tc.in)
		arr := b.NewArray()
		if tc.want == nil {
			require.True(t, arr.IsNull(0), "%s %v should be null", tc.kind, tc.in)
		} else {
			require.Equal(t, tc.want, cell(arr), "%s %v", tc.kind, tc.in)
		}
		arr.Release()
		b.Release()
	}
}

func TestAppendValueRejectsWrongShapes(t *testing.T) {
	cases := []struct {
		kind string
		in   any
	}{
		{"integer", "7"},
		{"char", json.Number("7")},
		{"date", "31/08/2026"},
		{"many2one", []any{json.Number("12")}},
	}
	for _, tc := range cases {
		cols, err := columnsFor([]string{"f"}, map[string]fieldInfo{"f": {tc.kind}})
		require.NoError(t, err)
		b := array.NewBuilder(memory.DefaultAllocator, cols[0].typ)
		require.Error(t, cols[0].appendValue(b, tc.in), "%s %v", tc.kind, tc.in)
		b.Release()
	}
}

// cell returns the first value of a one-element array as a plain Go value.
func cell(a arrow.Array) any {
	switch a := a.(type) {
	case *array.Int64:
		return a.Value(0)
	case *array.Float64:
		return a.Value(0)
	case *array.String:
		return a.Value(0)
	case *array.Boolean:
		return a.Value(0)
	case *array.Date32:
		return a.Value(0).ToTime().Format(time.DateOnly)
	case *array.Timestamp:
		return time.UnixMicro(int64(a.Value(0))).UTC().Format(time.DateTime)
	}
	return nil
}
