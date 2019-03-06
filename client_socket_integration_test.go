// +build integration

package tmi_test

import (
	"testing"

	"github.com/sunspots/tmi"
)

func TestTCPClient(t *testing.T) {
	testClient(t, tmi.NewTCPSocket(false))
}

func TestTCPSecureClient(t *testing.T) {
	testClient(t, tmi.NewTCPSocket(true))
}

func TestWSClient(t *testing.T) {
	testClient(t, tmi.NewWebSocket(false))
}

func TestWSSecureClient(t *testing.T) {
	testClient(t, tmi.NewWebSocket(true))
}
