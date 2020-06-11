package da

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewClient(t *testing.T) {
	client := New(testAppID)
	if client == nil {
		t.Fatal("shoud not be return nil.")
	}
}

const testAppID = "test_application_id"

const testText = "うちの庭には二羽鶏がいます。"

const testResultJSON = `{
  "id":"%s",
  "jsonrpc":"2.0",
  "result":{
    "chunks":[
      {
        "head":1,
        "id":0,
        "tokens":[
          ["うち","うち","うち","名詞","普通名詞","*","*"],
          ["の","の","の","助詞","接続助詞","*","*"]
        ]
      },
      {
        "head":3,
        "id":1,
        "tokens":[
          ["庭","にわ","庭","名詞","普通名詞","*","*"],
          ["に","に","に","助詞","格助詞","*","*"],
          ["は","は","は","助詞","副助詞","*","*"]
        ]
      },
      {
        "head":3,
        "id":2,
        "tokens":[
          ["二","に","二","名詞","数詞","*","*"],
          ["羽","わ","羽","接尾辞","名詞性名詞助数辞","*","*"],
          ["鶏","にわとり","鶏","名詞","普通名詞","*","*"],
          ["が","が","が","助詞","格助詞","*","*"]
        ]
      },
      {
        "head":-1,
        "id":3,
        "tokens":[
          ["い","い","いる","動詞","*","母音動詞","基本連用形"],
          ["ます","ます","ます","接尾辞","動詞性接尾辞","動詞性接尾辞ます型","基本形"]
        ]
      }
    ]
  }
}`

var testResult = Result{
	Chunks: []Chunk{
		{
			Head: 1,
			ID:   0,
			Tokens: []Token{
				Token{"うち", "うち", "うち", "名詞", "普通名詞", "*", "*"},
				Token{"の", "の", "の", "助詞", "接続助詞", "*", "*"},
			},
		},
		{
			Head: 3,
			ID:   1,
			Tokens: []Token{
				Token{"庭", "にわ", "庭", "名詞", "普通名詞", "*", "*"},
				Token{"に", "に", "に", "助詞", "格助詞", "*", "*"},
				Token{"は", "は", "は", "助詞", "副助詞", "*", "*"},
			},
		},
		{
			Head: 3,
			ID:   2,
			Tokens: []Token{
				Token{"二", "に", "二", "名詞", "数詞", "*", "*"},
				Token{"羽", "わ", "羽", "接尾辞", "名詞性名詞助数辞", "*", "*"},
				Token{"鶏", "にわとり", "鶏", "名詞", "普通名詞", "*", "*"},
				Token{"が", "が", "が", "助詞", "格助詞", "*", "*"},
			},
		},
		{
			Head: -1,
			ID:   3,
			Tokens: []Token{
				Token{"い", "い", "いる", "動詞", "*", "母音動詞", "基本連用形"},
				Token{"ます", "ます", "ます", "接尾辞", "動詞性接尾辞", "動詞性接尾辞ます型", "基本形"},
			},
		},
	},
}

type testReq struct {
	ID string `json:"id"`
}

func Test_Parse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const expectedUserAgent = "Yahoo AppID: " + testAppID
		userAgent := r.UserAgent()
		if userAgent != expectedUserAgent {
			t.Errorf("got %s, want %s", userAgent, expectedUserAgent)
		}
		var req testReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Error(err)
		}

		result := fmt.Sprintf(testResultJSON, req.ID)

		w.Write([]byte(result))
	}))

	APIEndpoint = ts.URL
	client := New(testAppID)
	res, err := client.Parse(context.Background(), testText)
	if err != nil {
		t.Fatalf("shoud not be fail: %v", err)
	}

	if diff := cmp.Diff(res, testResult); diff != "" {
		t.Errorf("Client.Parse, differs: (-got +want)\n%s", diff)
	}
}
