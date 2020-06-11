// Package jsonrpc is a JSON-RPC 2.0 client for Yahoo Japan Web API.
package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

const JSONRPCVersion = "2.0"

type Client struct {
	// endpoint is an endpoint of API.
	endpoint string
	// appID is an application ID for Yahoo Japan Web API.
	appID string

	// Client is a HTTP client. Use http.DefaultClient if it's nil.
	Client *http.Client
}

func NewClient(endpoint, appID string) *Client {
	return &Client{
		endpoint: endpoint,
		appID:    appID,
	}
}

type rpcRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      uuid.UUID   `json:"id"`
}

func newRPCRequest(method string, params interface{}) *rpcRequest {
	return &rpcRequest{
		JSONRPC: JSONRPCVersion,
		Method:  method,
		Params:  params,
		ID:      uuid.New(),
	}
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *ResponseError  `json:"error"`
	ID      uuid.UUID       `json:"id"`
}

type ErrorCode int

const (
	ParseError     ErrorCode = -32700
	InvalidRequest ErrorCode = -32600
	MethodNotFound ErrorCode = -32601
	InvalidParams  ErrorCode = -32602
	InternalError  ErrorCode = -32603

	ServerErrorStart ErrorCode = -32000
	ServerErrorEnd   ErrorCode = -32099
)

func IsServerError(code ErrorCode) bool {
	return code >= ServerErrorStart && code <= ServerErrorEnd
}

type ResponseError struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (err *ResponseError) Error() string {
	return fmt.Sprintf("%s (%d)", err.Message, err.Code)
}

func (client *Client) Call(ctx context.Context, method string, params interface{}, result interface{}) error {
	if method == "" {
		return errors.New("method is empty")
	}

	rpcReq := newRPCRequest(method, params)

	b, err := json.Marshal(rpcReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := client.newRequest(ctx, bytes.NewReader(b))
	if err != nil {
		return err
	}

	res, err := client.httpClient().Do(req)
	if err != nil {
		return fmt.Errorf("failed to post request: %w", err)
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server does not respond 200 OK: %s", res.Status)
	}

	var rpcRes rpcResponse

	if err := json.NewDecoder(res.Body).Decode(&rpcRes); err != nil {
		return fmt.Errorf("failed to decode response JSON: %w", err)
	}

	if rpcRes.Error != nil {
		return rpcRes.Error
	}

	if rpcRes.ID != rpcReq.ID {
		return errors.New("response ID is not matched to request")
	}

	if err := json.Unmarshal([]byte(rpcRes.Result), result); err != nil {
		return fmt.Errorf("failed to decode result JSON: %w", err)
	}

	return nil
}

func (client *Client) newRequest(ctx context.Context, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, client.endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new HTTP request: %w", err)
	}

	req.Header.Set("User-Agent", "Yahoo AppID: "+client.appID)

	return req, nil
}

func (client *Client) httpClient() *http.Client {
	if client.Client == nil {
		return http.DefaultClient
	}

	return client.Client
}
