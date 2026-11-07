// Command openapi converts the Swagger 2.0 document produced by `swag init`
// (docs/swagger.json) into an OpenAPI 3.0 document (docs/openapi.json,
// docs/openapi.yaml). swag itself only emits Swagger 2.0, so this is a
// second, pure-Go step — no Node.js/npx tooling needed — to get a proper
// OpenAPI 3 spec out of the Go annotations.
package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"gopkg.in/yaml.v3"
)

const (
	swaggerV2Path = "docs/swagger.json"
	openAPIJSON   = "docs/openapi.json"
	openAPIYAML   = "docs/openapi.yaml"
)

func main() {
	swaggerJSON, err := os.ReadFile(swaggerV2Path)
	if err != nil {
		log.Fatalf("failed to read %s (run `go tool swag init` first): %v", swaggerV2Path, err)
	}

	var doc2 openapi2.T
	if err := json.Unmarshal(swaggerJSON, &doc2); err != nil {
		log.Fatalf("failed to parse %s: %v", swaggerV2Path, err)
	}

	doc3, err := openapi2conv.ToV3(&doc2)
	if err != nil {
		log.Fatalf("failed to convert to OpenAPI 3: %v", err)
	}

	jsonBytes, err := json.MarshalIndent(doc3, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal OpenAPI 3 JSON: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(openAPIJSON), 0o750); err != nil {
		log.Fatalf("failed to create docs directory: %v", err)
	}
	if err := os.WriteFile(openAPIJSON, jsonBytes, 0o600); err != nil {
		log.Fatalf("failed to write %s: %v", openAPIJSON, err)
	}

	var generic any
	if err := json.Unmarshal(jsonBytes, &generic); err != nil {
		log.Fatalf("failed to prepare YAML output: %v", err)
	}
	yamlBytes, err := yaml.Marshal(generic)
	if err != nil {
		log.Fatalf("failed to marshal OpenAPI 3 YAML: %v", err)
	}
	if err := os.WriteFile(openAPIYAML, yamlBytes, 0o600); err != nil {
		log.Fatalf("failed to write %s: %v", openAPIYAML, err)
	}

	log.Printf("wrote %s and %s", openAPIJSON, openAPIYAML)
}
