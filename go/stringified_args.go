package main

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var toolArrayFields = map[string][]string{
	"submit_url_batch":    {"url_list"},
	"submit_url_indexnow": {"url_list"},
}

func coerceStringifiedArrayArgs(arrayFieldsByTool map[string][]string) mcp.Middleware {
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			call, ok := req.(*mcp.CallToolRequest)
			if !ok || method != "tools/call" {
				return next(ctx, method, req)
			}
			fields := arrayFieldsByTool[call.Params.Name]
			if len(fields) == 0 || len(call.Params.Arguments) == 0 {
				return next(ctx, method, req)
			}

			var args map[string]json.RawMessage
			if err := json.Unmarshal(call.Params.Arguments, &args); err != nil {
				return next(ctx, method, req)
			}

			changed := false
			for _, field := range fields {
				if coerced, ok := coerceStringifiedArray(args[field]); ok {
					args[field] = coerced
					changed = true
				}
			}
			if changed {
				if rewritten, err := json.Marshal(args); err == nil {
					call.Params.Arguments = rewritten
				}
			}
			return next(ctx, method, req)
		}
	}
}

func coerceStringifiedArray(raw json.RawMessage) (json.RawMessage, bool) {
	if len(raw) == 0 {
		return nil, false
	}
	var asString string
	if err := json.Unmarshal(raw, &asString); err != nil {
		return nil, false
	}
	var probe []json.RawMessage
	if err := json.Unmarshal([]byte(asString), &probe); err != nil {
		return nil, false
	}
	return json.RawMessage(asString), true
}
