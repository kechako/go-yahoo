// Package jsonrpc is a JSON-RPC 2.0 client for Yahoo Japan Web API.
package jsonrpc

import (
	"context"
	"errors"
	"net/http"

	"github.com/kechako/go-jsonrpc"
)

type Client struct {
	c        *jsonrpc.Client // JSON-RPC 2.0 client.
	endpoint string          // endpoint is an endpoint of API.
	appID    string          // appID is an application ID for Yahoo Japan Web API.
}

func NewClient(endpoint, appID string, opts ...Option) *Client {
	c := &Client{
		c:        &jsonrpc.Client{},
		endpoint: endpoint,
		appID:    appID,
	}

	if len(opts) > 0 {
		var copts clientOpts

		for _, opt := range opts {
			opt.apply(&copts)
		}

		c.c.HTTPClient = copts.httpClient
	}

	return c
}

func (client *Client) Call(ctx context.Context, method string, params interface{}, result interface{}) error {
	if method == "" {
		return errors.New("method is empty")
	}

	header := make(http.Header)
	header.Add("User-Agent", "Yahoo AppID: "+client.appID)

	err := client.c.Call(ctx, client.endpoint, method, params, result, jsonrpc.WithHeader(header))
	if err != nil {
		return err
	}

	return nil
}
