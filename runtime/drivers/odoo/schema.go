package odoo

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
)

// fieldInfo is what fields_get returns for a field when asked for the "type" attribute.
type fieldInfo struct {
	Type string `json:"type"`
}

// column is one output column. A many2one field becomes two: the id, named like the field,
// and the display name, with a "_name" suffix.
type column struct {
	name  string
	field string
	kind  string
	part  int // many2one only: 0 = id, 1 = display name
	typ   arrow.DataType
}

// Odoo sends datetimes as UTC strings without a zone.
var utcMicros = &arrow.TimestampType{Unit: arrow.Microsecond, TimeZone: "UTC"}

// columnsFor maps the requested fields, in order, to output columns. Types that do not hold one
// scalar per record (one2many, many2many, binary, json…) are rejected: the related model is
// read instead of flattening a list into a cell.
func columnsFor(fields []string, info map[string]fieldInfo) ([]column, error) {
	cols := make([]column, 0, len(fields))
	for _, f := range fields {
		fi, ok := info[f]
		if !ok {
			return nil, fmt.Errorf("odoo: the model has no field %q", f)
		}
		c := column{name: f, field: f, kind: fi.Type}
		switch fi.Type {
		case "integer":
			c.typ = arrow.PrimitiveTypes.Int64
		case "float", "monetary":
			c.typ = arrow.PrimitiveTypes.Float64
		case "char", "text", "html", "selection":
			c.typ = arrow.BinaryTypes.String
		case "boolean":
			c.typ = arrow.FixedWidthTypes.Boolean
		case "date":
			c.typ = arrow.FixedWidthTypes.Date32
		case "datetime":
			c.typ = utcMicros
		case "many2one":
			c.typ = arrow.PrimitiveTypes.Int64
			cols = append(cols, c, column{name: f + "_name", field: f, kind: fi.Type, part: 1, typ: arrow.BinaryTypes.String})
			continue
		default:
			return nil, fmt.Errorf("odoo: field %q has type %q, which the driver does not read; read the related model instead", f, fi.Type)
		}
		cols = append(cols, c)
	}
	return cols, nil
}

func schemaFor(cols []column) *arrow.Schema {
	fields := make([]arrow.Field, len(cols))
	for i, c := range cols {
		fields[i] = arrow.Field{Name: c.name, Type: c.typ, Nullable: true}
	}
	return arrow.NewSchema(fields, nil)
}

// appendValue writes one value of a search_read record into the column's builder. Odoo answers
// false for an empty field of any type, so false is null unless the field is a boolean.
func (c column) appendValue(b array.Builder, v any) error {
	if v == nil || (v == false && c.kind != "boolean") {
		b.AppendNull()
		return nil
	}
	switch c.kind {
	case "integer":
		n, err := asInt64(v)
		if err != nil {
			return err
		}
		b.(*array.Int64Builder).Append(n)
	case "float", "monetary":
		num, ok := v.(json.Number)
		if !ok {
			return fmt.Errorf("expected a number, got %T", v)
		}
		f, err := num.Float64()
		if err != nil {
			return err
		}
		b.(*array.Float64Builder).Append(f)
	case "char", "text", "html", "selection":
		s, ok := v.(string)
		if !ok {
			return fmt.Errorf("expected a string, got %T", v)
		}
		b.(*array.StringBuilder).Append(s)
	case "boolean":
		bv, ok := v.(bool)
		if !ok {
			return fmt.Errorf("expected a boolean, got %T", v)
		}
		b.(*array.BooleanBuilder).Append(bv)
	case "date":
		t, err := parseTime(time.DateOnly, v)
		if err != nil {
			return err
		}
		b.(*array.Date32Builder).Append(arrow.Date32FromTime(t))
	case "datetime":
		t, err := parseTime(time.DateTime, v)
		if err != nil {
			return err
		}
		b.(*array.TimestampBuilder).Append(arrow.Timestamp(t.UnixMicro()))
	case "many2one":
		pair, ok := v.([]any)
		if !ok || len(pair) != 2 {
			return fmt.Errorf("expected [id, name], got %v", v)
		}
		if c.part == 0 {
			n, err := asInt64(pair[0])
			if err != nil {
				return err
			}
			b.(*array.Int64Builder).Append(n)
			return nil
		}
		s, ok := pair[1].(string)
		if !ok {
			return fmt.Errorf("expected a display name, got %T", pair[1])
		}
		b.(*array.StringBuilder).Append(s)
	}
	return nil
}

func asInt64(v any) (int64, error) {
	num, ok := v.(json.Number)
	if !ok {
		return 0, fmt.Errorf("expected an integer, got %T", v)
	}
	return num.Int64()
}

func parseTime(layout string, v any) (time.Time, error) {
	s, ok := v.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("expected a %q string, got %T", layout, v)
	}
	return time.ParseInLocation(layout, s, time.UTC)
}
