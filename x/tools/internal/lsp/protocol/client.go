// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"context"
	"encoding/json"

	"github.com/Andyfoo/golang/x/tools/internal/jsonrpc2"
	"github.com/Andyfoo/golang/x/tools/internal/lsp/xlog"
)

type Client interface {
	ShowMessage(context.Context, *ShowMessageParams) error
	ShowMessageRequest(context.Context, *ShowMessageRequestParams) (*MessageActionItem, error)
	LogMessage(context.Context, *LogMessageParams) error
	Telemetry(context.Context, interface{}) error
	RegisterCapability(context.Context, *RegistrationParams) error
	UnregisterCapability(context.Context, *UnregistrationParams) error
	WorkspaceFolders(context.Context) ([]WorkspaceFolder, error)
	Configuration(context.Context, *ConfigurationParams) ([]interface{}, error)
	ApplyEdit(context.Context, *ApplyWorkspaceEditParams) (bool, error)
	PublishDiagnostics(context.Context, *PublishDiagnosticsParams) error
}

func clientHandler(log xlog.Logger, client Client) jsonrpc2.Handler {
	return func(ctx context.Context, conn *jsonrpc2.Conn, r *jsonrpc2.Request) {
		switch r.Method {
		case "$/cancelRequest":
			var params CancelParams
			if err := json.Unmarshal(*r.Params, &params); err != nil {
				sendParseError(ctx, log, conn, r, err)
				return
			}
			conn.Cancel(params.ID)

		case "window/showMessage":
			var params ShowMessageParams
			if err := json.Unmarshal(*r.Params, &params); err != nil {
				sendParseError(ctx, log, conn, r, err)
				return
			}
			if err := client.ShowMessage(ctx, &params); err != nil {
				log.Errorf(ctx, "%v", err)
			}

		case "window/showMessageRequest":
			var params ShowMessageRequestParams
			if err := json.Unmarshal(*r.Params, &params); err != nil {
				sendParseError(ctx, log, conn, r, err)
				return
			}
			resp, err := client.ShowMessageRequest(ctx, &params)
			if err := conn.Reply(ctx, r, resp, err); err != nil {
				log.Errorf(ctx, "%v", err)
			}

		case "window/logMessage":
			var params LogMessageParams
			if err := json.Unmarshal(*r.Params, &params); err != nil {
				sendParseError(ctx, log, conn, r, err)
				return
			}
			if err := client.LogMessage(ctx, &params); err != nil {
				log.Errorf(ctx, "%v", err)
			}

		case "telemetry/event":
			var params interface{}
			if err := json.Unmarshal(*r.Params, &params); err != nil {
				sendParseError(ctx, log, conn, r, err)
				return
			}
			if err := client.Telemetry(ctx, &params); err != nil {
				log.Errorf(ctx, "%v", err)
			}

		case "client/registerCapability":
			var params RegistrationParams
			if err := json.Unmarshal(*r.Params, &params); err != nil {
				sendParseError(ctx, log, conn, r, err)
				return
			}
			if err := client.RegisterCapability(ctx, &params); err != nil {
				log.Errorf(ctx, "%v", err)
			}

		case "client/unregisterCapability":
			var params UnregistrationParams
			if err := json.Unmarshal(*r.Params, &params); err != nil {
				sendParseError(ctx, log, conn, r, err)
				return
			}
			if err := client.UnregisterCapability(ctx, &params); err != nil {
				log.Errorf(ctx, "%v", err)
			}

		case "workspace/workspaceFolders":
			if r.Params != nil {
				conn.Reply(ctx, r, nil, jsonrpc2.NewErrorf(jsonrpc2.CodeInvalidParams, "Expected no params"))
				return
			}
			resp, err := client.WorkspaceFolders(ctx)
			if err := conn.Reply(ctx, r, resp, err); err != nil {
				log.Errorf(ctx, "%v", err)
			}

		case "workspace/configuration":
			var params ConfigurationParams
			if err := json.Unmarshal(*r.Params, &params); err != nil {
				sendParseError(ctx, log, conn, r, err)
				return
			}
			resp, err := client.Configuration(ctx, &params)
			if err := conn.Reply(ctx, r, resp, err); err != nil {
				log.Errorf(ctx, "%v", err)
			}

		case "workspace/applyEdit":
			var params ApplyWorkspaceEditParams
			if err := json.Unmarshal(*r.Params, &params); err != nil {
				sendParseError(ctx, log, conn, r, err)
				return
			}
			resp, err := client.ApplyEdit(ctx, &params)
			if err := conn.Reply(ctx, r, resp, err); err != nil {
				log.Errorf(ctx, "%v", err)
			}

		case "textDocument/publishDiagnostics":
			var params PublishDiagnosticsParams
			if err := json.Unmarshal(*r.Params, &params); err != nil {
				sendParseError(ctx, log, conn, r, err)
				return
			}
			if err := client.PublishDiagnostics(ctx, &params); err != nil {
				log.Errorf(ctx, "%v", err)
			}

		default:
			if r.IsNotify() {
				conn.Reply(ctx, r, nil, jsonrpc2.NewErrorf(jsonrpc2.CodeMethodNotFound, "method %q not found", r.Method))
			}
		}
	}
}

type clientDispatcher struct {
	*jsonrpc2.Conn
}

func (c *clientDispatcher) ShowMessage(ctx context.Context, params *ShowMessageParams) error {
	return c.Conn.Notify(ctx, "window/showMessage", params)
}

func (c *clientDispatcher) ShowMessageRequest(ctx context.Context, params *ShowMessageRequestParams) (*MessageActionItem, error) {
	var result MessageActionItem
	if err := c.Conn.Call(ctx, "window/showMessageRequest", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *clientDispatcher) LogMessage(ctx context.Context, params *LogMessageParams) error {
	return c.Conn.Notify(ctx, "window/logMessage", params)
}

func (c *clientDispatcher) Telemetry(ctx context.Context, params interface{}) error {
	return c.Conn.Notify(ctx, "telemetry/event", params)
}

func (c *clientDispatcher) RegisterCapability(ctx context.Context, params *RegistrationParams) error {
	return c.Conn.Notify(ctx, "client/registerCapability", params)
}

func (c *clientDispatcher) UnregisterCapability(ctx context.Context, params *UnregistrationParams) error {
	return c.Conn.Notify(ctx, "client/unregisterCapability", params)
}

func (c *clientDispatcher) WorkspaceFolders(ctx context.Context) ([]WorkspaceFolder, error) {
	var result []WorkspaceFolder
	if err := c.Conn.Call(ctx, "workspace/workspaceFolders", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *clientDispatcher) Configuration(ctx context.Context, params *ConfigurationParams) ([]interface{}, error) {
	var result []interface{}
	if err := c.Conn.Call(ctx, "workspace/configuration", params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *clientDispatcher) ApplyEdit(ctx context.Context, params *ApplyWorkspaceEditParams) (bool, error) {
	var result bool
	if err := c.Conn.Call(ctx, "workspace/applyEdit", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

func (c *clientDispatcher) PublishDiagnostics(ctx context.Context, params *PublishDiagnosticsParams) error {
	return c.Conn.Notify(ctx, "textDocument/publishDiagnostics", params)
}
