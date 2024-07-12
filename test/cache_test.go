package test

import (
	"CacheTest/cache"
	"testing"
)

func TestLRUCache(t *testing.T) {
	c := cache.NewLRUCache(2)

	c.Set("a", 1)
	c.Set("b", 2)

	if v, _ := c.Get("a"); v != 1 {
		t.Fatalf("Expected 1, got %v", v)
	}

	c.Set("c", 3)

	if _, err := c.Get("b"); err == nil {
		t.Fatalf("Expected error, got value")
	}

	if v, _ := c.Get("c"); v != 3 {
		t.Fatalf("Expected 3, got %v", v)
	}

	c.Remove("a")
	if _, err := c.Get("a"); err == nil {
		t.Fatalf("Expected error, got value")
	}

	c.Clear()
	if _, err := c.Get("c"); err == nil {
		t.Fatalf("Expected error, got value")
	}
}

func TestLFUCache(t *testing.T) {
	c := cache.NewLFUCache(2)
	c.Put(1, 1)
	c.Put(2, 2)

	if v, _ := c.Get(1); v != 1 {
		t.Fatalf("Expected 1, got %v", v)
	}
}

func TestTimeCache(t *testing.T) {
	c := cache.NewTimeCache(2)
	c.Put(1, 1)
	c.Put(2, 2)

	if v, _ := c.Get(1); v != 1 {
		t.Fatalf("Expected 1, got %v", v)
	}
}
