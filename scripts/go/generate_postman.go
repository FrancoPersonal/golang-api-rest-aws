//go:build ignore

// generate_postman.go reads docs/openapi.json (Swagger 2.0) and produces a
// Postman Collection v2.1 JSON at docs/Cauciones-API.postman_collection.json.
// Run via: go run ./scripts/generate_postman.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// ---------------------------------------------------------------------------
// Swagger 2.0 input types
// ---------------------------------------------------------------------------

type swaggerSpec struct {
	Info        swaggerInfo                       `json:"info"`
	Host        string                            `json:"host"`
	BasePath    string                            `json:"basePath"`
	Paths       map[string]map[string]swaggerOp   `json:"paths"`
	Definitions map[string]swaggerDefinition      `json:"definitions"`
}

type swaggerInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type swaggerOp struct {
	Summary    string          `json:"summary"`
	Tags       []string        `json:"tags"`
	Parameters []swaggerParam  `json:"parameters"`
	Consumes   []string        `json:"consumes"`
}

type swaggerParam struct {
	Name   string      `json:"name"`
	In     string      `json:"in"`
	Schema *swaggerRef `json:"schema"`
}

type swaggerRef struct {
	Ref string `json:"$ref"`
}

type swaggerDefinition struct {
	Properties map[string]swaggerProperty `json:"properties"`
}

type swaggerProperty struct {
	Type    string      `json:"type"`
	Example interface{} `json:"example"`
}

// ---------------------------------------------------------------------------
// Postman Collection v2.1 output types
// ---------------------------------------------------------------------------

type postmanCollection struct {
	Info     postmanInfo      `json:"info"`
	Item     []postmanFolder  `json:"item"`
	Variable []postmanVar     `json:"variable"`
}

type postmanInfo struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Schema      string `json:"schema"`
}

type postmanFolder struct {
	Name string           `json:"name"`
	Item []postmanRequest `json:"item"`
}

type postmanRequest struct {
	Name    string         `json:"name"`
	Request postmanReqBody `json:"request"`
}

type postmanReqBody struct {
	Method string          `json:"method"`
	Header []postmanHeader `json:"header"`
	Body   *postmanBody    `json:"body,omitempty"`
	URL    postmanURL      `json:"url"`
}

type postmanHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type postmanBody struct {
	Mode    string `json:"mode"`
	Raw     string `json:"raw"`
	Options struct {
		Raw struct {
			Language string `json:"language"`
		} `json:"raw"`
	} `json:"options"`
}

type postmanURL struct {
	Raw  string   `json:"raw"`
	Host []string `json:"host"`
	Path []string `json:"path"`
}

type postmanVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

func main() {
	const specPath = "docs/openapi.json"
	const outPath = "docs/Cauciones-API.postman_collection.json"

	raw, err := os.ReadFile(specPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot read %s: %v\n", specPath, err)
		os.Exit(1)
	}

	var spec swaggerSpec
	if err := json.Unmarshal(raw, &spec); err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot parse %s: %v\n", specPath, err)
		os.Exit(1)
	}

	collection := buildCollection(spec)

	out, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot marshal collection: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outPath, out, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot write %s: %v\n", outPath, err)
		os.Exit(1)
	}

	fmt.Printf("✓ Postman collection generated: %s\n", outPath)
}

// ---------------------------------------------------------------------------
// Builder
// ---------------------------------------------------------------------------

func buildCollection(spec swaggerSpec) postmanCollection {
	// Group operations by first tag to build folders.
	type opEntry struct {
		path   string
		method string
		op     swaggerOp
	}
	tagOps := map[string][]opEntry{}

	// Sort paths for deterministic output.
	paths := make([]string, 0, len(spec.Paths))
	for p := range spec.Paths {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	methodOrder := []string{"get", "post", "put", "patch", "delete"}

	for _, path := range paths {
		methods := spec.Paths[path]
		for _, method := range methodOrder {
			op, ok := methods[method]
			if !ok {
				continue
			}
			tag := "General"
			if len(op.Tags) > 0 {
				tag = op.Tags[0]
			}
			tagOps[tag] = append(tagOps[tag], opEntry{path, method, op})
		}
	}

	// Sort tags for deterministic folder order.
	tags := make([]string, 0, len(tagOps))
	for t := range tagOps {
		tags = append(tags, t)
	}
	sort.Strings(tags)

	folders := make([]postmanFolder, 0, len(tags))
	for _, tag := range tags {
		folder := postmanFolder{Name: tag}
		for _, e := range tagOps[tag] {
			folder.Item = append(folder.Item, buildRequest(e.path, e.method, e.op, spec))
		}
		folders = append(folders, folder)
	}

	base := strings.TrimSuffix(spec.Host, "/")
	if spec.BasePath != "/" && spec.BasePath != "" {
		base += spec.BasePath
	}

	return postmanCollection{
		Info: postmanInfo{
			Name:        spec.Info.Title,
			Description: spec.Info.Description,
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Item: folders,
		Variable: []postmanVar{
			{Key: "base_url", Value: "http://" + base, Type: "string"},
			{Key: "jwt_token", Value: "your-jwt-token-here", Type: "string"},
			{Key: "id", Value: "550e8400-e29b-41d4-a716-446655440000", Type: "string"},
		},
	}
}

func buildRequest(path, method string, op swaggerOp, spec swaggerSpec) postmanRequest {
	name := op.Summary
	if name == "" {
		name = strings.ToUpper(method) + " " + path
	}

	headers := []postmanHeader{
		{Key: "Authorization", Value: "Bearer {{jwt_token}}"},
	}

	// Add Content-Type for operations that consume JSON.
	for _, p := range op.Parameters {
		if p.In == "body" {
			headers = append(headers, postmanHeader{Key: "Content-Type", Value: "application/json"})
			break
		}
	}

	var body *postmanBody
	for _, p := range op.Parameters {
		if p.In == "body" && p.Schema != nil {
			example := buildBodyExample(p.Schema.Ref, spec.Definitions)
			if example != "" {
				b := &postmanBody{Raw: example}
				b.Options.Raw.Language = "json"
				b.Mode = "raw"
				body = b
			}
			break
		}
	}

	// Replace path params with Postman variables (e.g. {id} → {{id}}).
	rawURL := "{{base_url}}" + strings.NewReplacer("{", "{{", "}", "}}").Replace(path)
	pathSegments := pathSegments(path)

	return postmanRequest{
		Name: name,
		Request: postmanReqBody{
			Method: strings.ToUpper(method),
			Header: headers,
			Body:   body,
			URL: postmanURL{
				Raw:  rawURL,
				Host: []string{"{{base_url}}"},
				Path: pathSegments,
			},
		},
	}
}

// buildBodyExample resolves a $ref to a definition and builds a JSON example
// object from the `example` values of each property.
func buildBodyExample(ref string, defs map[string]swaggerDefinition) string {
	// $ref format: "#/definitions/dto.CreateSuretyBondRequest"
	name := strings.TrimPrefix(ref, "#/definitions/")
	def, ok := defs[name]
	if !ok {
		return ""
	}

	obj := map[string]interface{}{}
	// Sort properties for deterministic output.
	props := make([]string, 0, len(def.Properties))
	for k := range def.Properties {
		props = append(props, k)
	}
	sort.Strings(props)

	for _, k := range props {
		v := def.Properties[k]
		if v.Example != nil {
			obj[k] = v.Example
		} else {
			obj[k] = placeholderFor(v.Type)
		}
	}

	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}

func placeholderFor(t string) interface{} {
	switch t {
	case "integer", "number":
		return 0
	case "boolean":
		return false
	case "array":
		return []interface{}{}
	default:
		return ""
	}
}

// pathSegments splits a path into Postman URL segments, converting {param} to {{param}}.
func pathSegments(path string) []string {
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	for i, p := range parts {
		if strings.HasPrefix(p, "{") && strings.HasSuffix(p, "}") {
			parts[i] = "{{" + p[1:len(p)-1] + "}}"
		}
	}
	return parts
}
