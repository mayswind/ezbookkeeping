package core

import "encoding/json"

// JSONRPCVersion defines the version of JSON-RPC protocol
const JSONRPCVersion = "2.0"

// JSONRPCRequest represents the JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      any             `json:"id,omitempty"`
}

// JSONRPCResponse represents the JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	Result  any           `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
	ID      any           `json:"id,omitempty"`
}

// JSONRPCError represents the JSON-RPC 2.0 error object
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// JSONRPCParseError represents the "Parse error" in JSON-RPC 2.0
var JSONRPCParseError = &JSONRPCError{
	Code:    -32700,
	Message: "Parse error",
	Data:    nil,
}

// JSONRPCMethodNotFoundError represents the "Method not found" error in JSON-RPC 2.0
var JSONRPCMethodNotFoundError = &JSONRPCError{
	Code:    -32601,
	Message: "Method not found",
	Data:    nil,
}

// JSONRPCInvalidParamsError represents the "Invalid params" error in JSON-RPC 2.0
var JSONRPCInvalidParamsError = &JSONRPCError{
	Code:    -32602,
	Message: "Invalid params",
	Data:    nil,
}

// JSONRPCInternalError represents the "Internal error" in JSON-RPC 2.0
var JSONRPCInternalError = &JSONRPCError{
	Code:    -32603,
	Message: "Internal error",
	Data:    nil,
}

// NewJSONRPCResponse creates a new JSON-RPC response with the result
func NewJSONRPCResponse(id any, result any) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		Result:  result,
		Error:   nil,
		ID:      id,
	}
}

// NewJSONRPCErrorResponse creates a new JSON-RPC error response
func NewJSONRPCErrorResponse(id any, err *JSONRPCError) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		Result:  nil,
		Error: &JSONRPCError{
			Code:    err.Code,
			Message: err.Message,
			Data:    nil,
		},
		ID: id,
	}
}

// NewJSONRPCErrorResponseWithCause creates a new JSON-RPC error response
func NewJSONRPCErrorResponseWithCause(id any, err *JSONRPCError, cause string) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		Result:  nil,
		Error: &JSONRPCError{
			Code:    err.Code,
			Message: err.Message,
			Data:    cause,
		},
		ID: id,
	}
}
