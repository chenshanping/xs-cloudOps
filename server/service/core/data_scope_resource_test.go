package core

import (
	"reflect"
	"testing"
)

func TestSupportedDataScopeResources(t *testing.T) {
	t.Run("resource metadata is complete", func(t *testing.T) {
		resourceType := reflect.TypeOf(DataScopeResource{})
		expectedFields := []struct {
			name    string
			jsonTag string
			kind    reflect.Kind
		}{
			{name: "Code", jsonTag: "code", kind: reflect.String},
			{name: "Label", jsonTag: "label", kind: reflect.String},
			{name: "Description", jsonTag: "description", kind: reflect.String},
			{name: "OwnerFields", jsonTag: "owner_fields", kind: reflect.Slice},
		}

		for _, expectedField := range expectedFields {
			field, ok := resourceType.FieldByName(expectedField.name)
			if !ok {
				t.Fatalf("expected DataScopeResource.%s to be declared", expectedField.name)
			}
			if field.Type.Kind() != expectedField.kind {
				t.Fatalf(
					"expected DataScopeResource.%s kind %s, got %s",
					expectedField.name,
					expectedField.kind,
					field.Type.Kind(),
				)
			}
			if tag := field.Tag.Get("json"); tag != expectedField.jsonTag {
				t.Fatalf(
					"expected DataScopeResource.%s json tag %q, got %q",
					expectedField.name,
					expectedField.jsonTag,
					tag,
				)
			}
		}

		resources := SupportedDataScopeResources()
		if len(resources) == 0 {
			t.Fatal("expected supported data scope resources")
		}

		resourcesByCode := indexResourcesByCode(resources)
		for _, code := range []string{
			DataScopeResourceUserManagement,
			DataScopeResourceDeptManagement,
		} {
			resourceValue, ok := resourcesByCode[code]
			if !ok {
				t.Fatalf("expected resource %s to be registered", code)
			}

			if label := resourceValue.FieldByName("Label").String(); label == "" {
				t.Fatalf("expected resource %s label to be non-empty", code)
			}
			descriptionField := resourceValue.FieldByName("Description")
			if !descriptionField.IsValid() || descriptionField.String() == "" {
				t.Fatalf("expected resource %s description to be non-empty", code)
			}

			ownerFields := extractStringSliceField(t, resourceValue, "OwnerFields")
			if len(ownerFields) == 0 {
				t.Fatalf("expected resource %s owner fields to be non-empty", code)
			}
		}
	})

	t.Run("resource contract is correct per resource", func(t *testing.T) {
		resources := SupportedDataScopeResources()
		if len(resources) == 0 {
			t.Fatal("expected supported data scope resources")
		}

		resourcesByCode := indexResourcesByCode(resources)
		assertResourceContract(t, resourcesByCode, DataScopeResourceUserManagement, "用户管理数据权限资源，支持按部门和创建人限定访问范围。", []string{DataScopeOwnerFieldDeptID, DataScopeOwnerFieldCreatedBy})
		assertResourceContract(t, resourcesByCode, DataScopeResourceDeptManagement, "部门管理数据权限资源，支持按部门范围限定可管理部门。", []string{DataScopeOwnerFieldDeptID})
	})

	t.Run("returns defensive copy", func(t *testing.T) {
		resources := SupportedDataScopeResources()
		if len(resources) == 0 {
			t.Fatal("expected supported data scope resources")
		}

		originalCode := resources[0].Code
		originalOwnerFields := extractStringSliceField(t, reflect.ValueOf(resources[0]), "OwnerFields")
		if len(originalOwnerFields) == 0 {
			t.Fatal("expected owner fields on first resource")
		}

		resources[0].Code = "mutated-code"
		ownerFieldsField := reflect.ValueOf(&resources[0]).Elem().FieldByName("OwnerFields")
		if !ownerFieldsField.IsValid() {
			t.Fatal("expected DataScopeResource.OwnerFields to be declared")
		}
		ownerFieldsField.Index(0).SetString("mutated-owner-field")

		freshResources := SupportedDataScopeResources()
		if len(freshResources) == 0 {
			t.Fatal("expected supported data scope resources")
		}

		if freshResources[0].Code != originalCode {
			t.Fatalf("expected fresh resource code %q, got %q", originalCode, freshResources[0].Code)
		}

		freshOwnerFields := extractStringSliceField(t, reflect.ValueOf(freshResources[0]), "OwnerFields")
		if len(freshOwnerFields) == 0 {
			t.Fatal("expected fresh owner fields on first resource")
		}
		if freshOwnerFields[0] != originalOwnerFields[0] {
			t.Fatalf("expected fresh owner field %q, got %q", originalOwnerFields[0], freshOwnerFields[0])
		}
	})
}

func indexResourcesByCode(resources []DataScopeResource) map[string]reflect.Value {
	indexed := make(map[string]reflect.Value, len(resources))
	for _, resource := range resources {
		resourceValue := reflect.ValueOf(resource)
		indexed[resource.Code] = resourceValue
	}
	return indexed
}

func assertResourceContract(
	t *testing.T,
	resourcesByCode map[string]reflect.Value,
	resourceCode string,
	expectedDescription string,
	expectedOwnerFields []string,
) {
	t.Helper()

	resourceValue, ok := resourcesByCode[resourceCode]
	if !ok {
		t.Fatalf("expected resource %s to be registered", resourceCode)
	}

	description := resourceValue.FieldByName("Description")
	if !description.IsValid() {
		t.Fatalf("expected resource %s description to be declared", resourceCode)
	}
	if description.String() != expectedDescription {
		t.Fatalf(
			"expected resource %s description %q, got %q",
			resourceCode,
			expectedDescription,
			description.String(),
		)
	}

	ownerFields := extractStringSliceField(t, resourceValue, "OwnerFields")
	if !reflect.DeepEqual(ownerFields, expectedOwnerFields) {
		t.Fatalf(
			"expected resource %s owner fields %v, got %v",
			resourceCode,
			expectedOwnerFields,
			ownerFields,
		)
	}
}

func extractStringSliceField(t *testing.T, resourceValue reflect.Value, fieldName string) []string {
	t.Helper()

	field := resourceValue.FieldByName(fieldName)
	if !field.IsValid() {
		t.Fatalf("expected DataScopeResource.%s to be declared", fieldName)
	}
	if field.Kind() != reflect.Slice {
		t.Fatalf("expected DataScopeResource.%s to be a slice, got %s", fieldName, field.Kind())
	}

	values := make([]string, field.Len())
	for i := 0; i < field.Len(); i++ {
		values[i] = field.Index(i).String()
	}
	return values
}
