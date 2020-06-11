package jsonrpc

import (
	"net/http"
	"testing"
)

const testAppID = "test_application_id"

func TestNewClient(t *testing.T) {
	const testEndpoint = "https://example.com/endpoint"
	client := NewClient(testEndpoint, testAppID)
	if client == nil {
		t.Fatal("must not return nil.")
	}

	if client.endpoint != testEndpoint {
		t.Errorf("Client.endpoint: got %s, want %s", client.endpoint, testEndpoint)
	}

	if client.appID != testAppID {
		t.Errorf("Client.appID: got %s, want %s", client.appID, testAppID)
	}

	if client.httpClient() == nil {
		t.Error("Client.httpClient() must not be nil")
	}

	if client.httpClient() != http.DefaultClient {
		t.Error("Client.httpClient() must be http.DefaultClient")
	}

	var httpClient = &http.Client{}
	client.Client = httpClient
	if client.httpClient() == nil {
		t.Error("Client.httpClient() must not be nil")
	}

	if client.httpClient() != httpClient {
		t.Error("Client.httpClient() must not be http.DefaultClient")
	}
}
