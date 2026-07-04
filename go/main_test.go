package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ncosentino/bing-webmaster-mcp/go/internal/bingwebmaster"
	"github.com/ncosentino/bing-webmaster-mcp/go/internal/indexnow"
)

func TestNewServer_RegistersTools(_ *testing.T) {
	srv := mcp.NewServer(&mcp.Implementation{
		Name:    "bing-webmaster-mcp",
		Version: "test",
	}, nil)

	var bingClient *bingwebmaster.Client
	var indexNowClient *indexnow.Client

	registerTools(srv, bingClient, indexNowClient)
}

func TestListSites_InputSchema(t *testing.T) {
	var schema map[string]any
	if err := json.Unmarshal(listSitesInputSchema, &schema); err != nil {
		t.Fatalf("unmarshal listSitesInputSchema: %v", err)
	}

	properties, ok := schema["properties"]
	if !ok {
		t.Error("list_sites InputSchema is missing the 'properties' field; strict MCP clients will reject it")
	} else if _, ok := properties.(map[string]any); !ok {
		t.Error("list_sites InputSchema 'properties' field must be a JSON object")
	}

	required, ok := schema["required"]
	if !ok {
		t.Error("list_sites InputSchema is missing the 'required' field")
	} else if _, ok := required.([]any); !ok {
		t.Error("list_sites InputSchema 'required' field must be a JSON array")
	}

	additionalProperties, ok := schema["additionalProperties"]
	if !ok {
		t.Error("list_sites InputSchema is missing the 'additionalProperties' field")
	} else if additionalPropertiesBool, ok := additionalProperties.(bool); !ok {
		t.Error("list_sites InputSchema 'additionalProperties' field must be a boolean")
	} else if additionalPropertiesBool {
		t.Error("list_sites InputSchema 'additionalProperties' must be false")
	}
}

func TestToolHandler_ReturnsErrorJSON(t *testing.T) {
	handler := toolHandler("testing", func(context.Context, struct{}) (any, error) {
		return nil, context.DeadlineExceeded
	})

	result, _, err := handler(context.Background(), nil, struct{}{})
	if err != nil {
		t.Fatalf("handler err = %v", err)
	}
	if len(result.Content) != 1 {
		t.Fatalf("len(Content) = %d, want 1", len(result.Content))
	}
}
