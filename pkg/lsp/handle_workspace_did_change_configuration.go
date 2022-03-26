package lsp

import (
	"context"
	"encoding/json"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) handleWorkspaceDidChangeConfiguration(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params DidChangeConfigurationParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	// TODO
	return nil, nil
}
