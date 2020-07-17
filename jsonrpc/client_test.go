package jsonrpc

import (
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
}
