package main_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	lookout "github.com/jakealves/atlas-game-lookout"
)

func TestPrintWebhook(t *testing.T) {
	var str bytes.Buffer
	log.SetOutput(&str)
	discordWebookJSON := map[string]interface{}{
		"content": "**Eastern Order (1672657922)**\nDay 28255, 04:01:05: Crew member Lyon Lint - Lvl 34 was killed by an Alpha Elephant - Lvl 223!\n",
	}
	body, _ := json.Marshal(discordWebookJSON)
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	err := lookout.PrintWebhook(req)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
	if str.String() == "" {
		t.Errorf("Expected log to be empty %v", str.String())
	}
}

func TestRelayWebhook(t *testing.T) {
	whs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "POST" {
			t.Errorf("Expected 'POST' request, got '%v'", r.Method)
		}
		if r.URL.EscapedPath() != "/webhook" {
			t.Errorf("Expected request to '/webhook', got '%v'", r.URL.EscapedPath())
		}
	}))
	defer whs.Close()
	req := httptest.NewRequest(http.MethodPost, "/upper?word=abc", nil)
	webhookURL := whs.URL + "/webhook"
	err := lookout.RelayWebhook(req, webhookURL)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
}

func TestExtractEvents(t *testing.T) {
	lookout.ExtractEvents("**Eastern Order (1672657922)**\nDay 28537, 01:46:42: Crew member Liv Kern - Lvl 40 was killed by TheOsky In Bike - Lvl 73 (Los Cheetos)!\nDay 28537, 01:50:51: Crew member Alfred Kern - Lvl 47 was killed by TheOsky In Bike - Lvl 73 (Los Cheetos)!\nDay 28537, 01:18:16: Alfred Kern claimed 'Crab - Lvl 45 (Crab)'!\n")
}
