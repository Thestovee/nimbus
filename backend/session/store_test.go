package session

import (
	"encoding/base64"
	"testing"
	"time"

	"nimbus/client"
)

func TestStoreCreate(t *testing.T) {
	store := NewStore()
	if _, err := store.Create(nil); err == nil {
		t.Fatal("Create(nil) should fail")
	}

	firstClient := client.NewClient()
	before := time.Now()
	firstToken, err := store.Create(firstClient)
	if err != nil {
		t.Fatalf("Create returned an error: %v", err)
	}
	after := time.Now()
	bytes, err := base64.RawURLEncoding.DecodeString(firstToken)
	if err != nil || len(bytes) < 32 {
		t.Fatalf("token is not at least 32 random bytes encoded as URL-safe base64")
	}
	firstSession, ok := store.sessions[firstToken]
	if !ok || firstSession.Client != firstClient {
		t.Fatal("created session does not contain its client")
	}
	if firstSession.ExpiresAt.Before(before.Add(24*time.Hour)) || firstSession.ExpiresAt.After(after.Add(24*time.Hour)) {
		t.Errorf("unexpected expiry: %v", firstSession.ExpiresAt)
	}

	secondClient := client.NewClient()
	secondToken, err := store.Create(secondClient)
	if err != nil {
		t.Fatalf("second Create returned an error: %v", err)
	}
	if secondToken == firstToken || store.sessions[secondToken].Client != secondClient {
		t.Fatal("two clients did not receive separate sessions")
	}
}

func TestStoreGet(t *testing.T) {
	store := NewStore()
	c := client.NewClient()
	token, err := store.Create(c)
	if err != nil {
		t.Fatalf("Create returned an error: %v", err)
	}

	got, ok := store.Get(token)
	if !ok || got != c {
		t.Fatal("Get did not return the client for a valid token")
	}
	if got, ok := store.Get("missing"); ok || got != nil {
		t.Fatal("Get accepted an unknown token")
	}

	store.sessions[token] = Session{Client: c, ExpiresAt: time.Now().Add(-time.Second)}
	if got, ok := store.Get(token); ok || got != nil {
		t.Fatal("Get accepted an expired token")
	}
	if _, exists := store.sessions[token]; exists {
		t.Fatal("Get did not remove the expired session")
	}
}
