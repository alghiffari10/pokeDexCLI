package pokecache

import (
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	cache := NewCache(5 * time.Second)

	key := "test"
	value := []byte("hello")

	cache.Add(key, value)

	got, ok := cache.Get(key)

	if !ok {
		t.Fatal("expected key to exist")
	}

	if string(got) != string(value) {
		t.Fatalf("expected %s, got %s", value, got)
	}
}
