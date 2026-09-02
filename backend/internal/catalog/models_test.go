package catalog

import "testing"

func TestCatalogHasCurrentImageAndVideoFamilies(t *testing.T) {
	catalog := BuildCatalog()
	if len(catalog.ByCategory(CatImage)) < 10 {
		t.Fatal("image catalog is unexpectedly small")
	}
	if len(catalog.ByCategory(CatVideo)) < 10 {
		t.Fatal("video catalog is unexpectedly small")
	}
	for _, id := range []string{"qwen-image-3.0-pro", "wan2.7-image-pro", "wan3.0-video-prime", "wan2.7-videoedit"} {
		if catalog.Find(id) == nil {
			t.Fatalf("catalog is missing %q", id)
		}
	}
}

func TestCatalogFieldsAreRoutable(t *testing.T) {
	for _, model := range BuildCatalog().Models {
		if model.ID == "" || model.Name == "" || model.Endpoint == "" || model.Payload == "" {
			t.Fatalf("model has incomplete routing metadata: %+v", model)
		}
		if model.Category != CatImage && model.Category != CatVideo {
			t.Fatalf("model %q has invalid category %q", model.ID, model.Category)
		}
		seen := map[string]bool{}
		for _, field := range model.Parameters {
			if field.Name == "" || field.Label == "" || field.Type == "" {
				t.Fatalf("model %q has incomplete field: %+v", model.ID, field)
			}
			if field.Scope != ScopeInput && field.Scope != ScopeParameters {
				t.Fatalf("model %q field %q has invalid scope %q", model.ID, field.Name, field.Scope)
			}
			if seen[field.Name] {
				t.Fatalf("model %q repeats field %q", model.ID, field.Name)
			}
			seen[field.Name] = true
		}
	}
}
