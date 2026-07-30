package powerdns

import (
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// The OpenAPI specification is vendored from the PowerDNS Authoritative Server
// repository and converted to JSON so it can be parsed with the standard
// library (no YAML dependency).
//
// Source: https://github.com/PowerDNS/pdns/blob/master/docs/http-api/openapi/authoritative-api-openapi.yaml
// Refresh: yq -o=json '.' authoritative-api-openapi.yaml > testdata/authoritative-api-openapi.json
const openAPISpecPath = "testdata/authoritative-api-openapi.json"

type openAPISpec struct {
	Components struct {
		Schemas map[string]struct {
			Properties map[string]json.RawMessage `json:"properties"`
		} `json:"schemas"`
	} `json:"components"`
}

func loadOpenAPISpec(t *testing.T) openAPISpec {
	t.Helper()

	raw, err := os.ReadFile(openAPISpecPath)
	if err != nil {
		t.Fatalf("read OpenAPI spec: %v", err)
	}

	var spec openAPISpec
	if err := json.Unmarshal(raw, &spec); err != nil {
		t.Fatalf("parse OpenAPI spec: %v", err)
	}

	return spec
}

// jsonFieldNames returns the JSON field names of a struct type, skipping fields
// tagged with "-".
func jsonFieldNames(t reflect.Type) map[string]struct{} {
	names := make(map[string]struct{})
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		name := strings.Split(tag, ",")[0]
		if name == "" || name == "-" {
			continue
		}
		names[name] = struct{}{}
	}
	return names
}

func sorted(set map[string]struct{}) []string {
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// specConformance describes how a Go struct maps to an OpenAPI schema.
type specConformance struct {
	schema string
	typ    reflect.Type
	// complete asserts the mapping is bidirectional: every spec property must
	// also exist on the Go struct. Set to false for structs that intentionally
	// implement only a subset of the schema.
	complete bool
	// ignore lists Go JSON field names that legitimately have no spec property
	// (e.g. undocumented but supported request fields).
	ignore []string
}

func TestStructsMatchOpenAPISpec(t *testing.T) {
	spec := loadOpenAPISpec(t)

	cases := []specConformance{
		{schema: "Cryptokey", typ: reflect.TypeOf(Cryptokey{}), complete: true},
		{schema: "TSIGKey", typ: reflect.TypeOf(TSIGKey{}), complete: true},
		{schema: "Metadata", typ: reflect.TypeOf(Metadata{}), complete: true},
		{schema: "Server", typ: reflect.TypeOf(Server{}), complete: true},
		{schema: "CacheFlushResult", typ: reflect.TypeOf(CacheFlushResult{}), complete: true},
		{schema: "ConfigSetting", typ: reflect.TypeOf(ConfigSetting{}), complete: true},
		{schema: "Comment", typ: reflect.TypeOf(Comment{}), complete: true},
		{schema: "Zone", typ: reflect.TypeOf(Zone{}), complete: false},
		{schema: "RRSet", typ: reflect.TypeOf(RRset{}), complete: false},
		{schema: "Record", typ: reflect.TypeOf(Record{}), complete: false, ignore: []string{"set-ptr"}},
		{schema: "Error", typ: reflect.TypeOf(Error{}), complete: false},
	}

	for _, tc := range cases {
		t.Run(tc.schema, func(t *testing.T) {
			schema, ok := spec.Components.Schemas[tc.schema]
			if !ok {
				t.Fatalf("schema %q not found in OpenAPI spec", tc.schema)
			}

			specProps := make(map[string]struct{}, len(schema.Properties))
			for name := range schema.Properties {
				specProps[name] = struct{}{}
			}

			goFields := jsonFieldNames(tc.typ)

			ignored := make(map[string]struct{}, len(tc.ignore))
			for _, name := range tc.ignore {
				ignored[name] = struct{}{}
			}

			for field := range goFields {
				if _, ok := ignored[field]; ok {
					continue
				}
				if _, ok := specProps[field]; !ok {
					t.Errorf("struct field %q (json) has no property in OpenAPI schema %q", field, tc.schema)
				}
			}

			if tc.complete {
				for prop := range specProps {
					if _, ok := goFields[prop]; !ok {
						t.Errorf("OpenAPI schema %q property %q is not implemented by the struct", tc.schema, prop)
					}
				}
			}
		})
	}

	t.Logf("validated %d structs against OpenAPI spec", len(cases))
}

// TestOpenAPISpecIsParseable guards against a corrupt or truncated vendored spec.
func TestOpenAPISpecIsParseable(t *testing.T) {
	spec := loadOpenAPISpec(t)

	if len(spec.Components.Schemas) == 0 {
		t.Fatal("OpenAPI spec contains no schemas")
	}

	if _, ok := spec.Components.Schemas["Cryptokey"]; !ok {
		t.Errorf("expected Cryptokey schema in spec, got: %v", sorted(schemaNames(spec)))
	}
}

func schemaNames(spec openAPISpec) map[string]struct{} {
	names := make(map[string]struct{}, len(spec.Components.Schemas))
	for name := range spec.Components.Schemas {
		names[name] = struct{}{}
	}
	return names
}
