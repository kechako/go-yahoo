// Package da is a client for Japanese dependency analysis API V2 served by Yahoo Japan.
package da

import (
	"context"
	"fmt"

	"github.com/kechako/go-yahoo/jsonrpc"
)

// APIEndpoint is an endpoint URL of API.
var APIEndpoint = "https://jlp.yahooapis.jp/DAService/V2/parse"

// A Client represents a client for API.
type Client struct {
	client *jsonrpc.Client
}

// New returns a new *Client that has appID.
func New(appID string, opts ...Option) *Client {
	var copts clientOpts

	for _, opt := range opts {
		opt.apply(&copts)
	}

	rpcClient := jsonrpc.NewClient(APIEndpoint, appID)
	rpcClient.Client = copts.httpClient

	return &Client{
		client: rpcClient,
	}
}

const ParseMethod = "jlp.daservice.parse"

// Parse returns a Result parsed japanese dependency.
func (c *Client) Parse(ctx context.Context, text string) (result Result, err error) {
	err = c.client.Call(ctx, ParseMethod, request{Text: text}, &result)
	if err != nil {
		err = fmt.Errorf("failed to parse text: %w", err)
		return
	}

	return
}
