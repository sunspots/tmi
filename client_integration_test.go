// +build integration

package tmi_test

import (
	"testing"

	"github.com/sunspots/tmi"
)

// See client_tcp_test.go and client_ws_test.go
func testClient(t *testing.T, socket tmi.Socket) {
	client := tmi.New("justinfan999", tmi.AnonymousUser, socket)

	err := client.Connect()
	if err != nil {
		t.Error(err)
	}
	m, err := client.ReadMessage()
	if err != nil {
		t.Error(err)
	}
	if m.From != "tmi.twitch.tv" {
		t.Errorf("Expected: '%v', Got: '%v'", "tmi.twitch.tv", m.From)
	}
}
