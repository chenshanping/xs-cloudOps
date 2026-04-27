package registry

import "testing"

type fieldsEmbeddedPage struct {
	Page     int `form:"page" comment:"页码"`
	PageSize int `form:"page_size" comment:"每页数量" binding:"required"`
}

type fieldsFixture struct {
	fieldsEmbeddedPage
	Keyword string `form:"keyword" comment:"关键字"`
	Status  int    `json:"status" comment:"状态"`
}

func TestParseDefinitionFieldsSharesTagAndRequiredRules(t *testing.T) {
	fields := ParseDefinitionFields(fieldsFixture{})
	if len(fields) != 4 {
		t.Fatalf("expected 4 definition fields, got %d", len(fields))
	}

	expected := map[string]struct {
		required bool
		typ      string
	}{
		"page":      {required: false, typ: "integer"},
		"page_size": {required: true, typ: "integer"},
		"keyword":   {required: false, typ: "string"},
		"status":    {required: false, typ: "integer"},
	}

	for _, field := range fields {
		want, ok := expected[field.Name]
		if !ok {
			t.Fatalf("unexpected field parsed: %s", field.Name)
		}
		if field.Required != want.required {
			t.Fatalf("field %s required mismatch: got %v want %v", field.Name, field.Required, want.required)
		}
		if field.SwaggerType != want.typ {
			t.Fatalf("field %s swagger type mismatch: got %s want %s", field.Name, field.SwaggerType, want.typ)
		}
	}
}

func TestParseStructFieldsUsesSharedMetadataForApiSync(t *testing.T) {
	fields := ParseStructFields(fieldsFixture{}, "query")
	if len(fields) != 4 {
		t.Fatalf("expected 4 api-sync fields, got %d", len(fields))
	}

	if fields[0].Name != "page" || fields[0].In != "query" {
		t.Fatalf("unexpected first field: %#v", fields[0])
	}
	if fields[1].Name != "page_size" || !fields[1].Required {
		t.Fatalf("unexpected second field: %#v", fields[1])
	}
	if fields[2].Name != "keyword" || fields[2].Description != "关键字" {
		t.Fatalf("unexpected third field: %#v", fields[2])
	}
	if fields[3].Name != "status" || fields[3].Type != "integer" {
		t.Fatalf("unexpected fourth field: %#v", fields[3])
	}
}
