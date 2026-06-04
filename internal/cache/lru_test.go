package cache

import "testing"

func TestLRUGetPut(t *testing.T) {
	c := New(2)

	c.Put("a", "1")
	c.Put("b", "2")

	if v, ok := c.Get("a"); !ok || v != "1" {
		t.Fatalf("expected a=1, got %q ok=%v", v, ok)
	}

	c.Put("c", "3") // evicts b

	if _, ok := c.Get("b"); ok {
		t.Fatal("expected b to be evicted")
	}
	if v, ok := c.Get("c"); !ok || v != "3" {
		t.Fatalf("expected c=3, got %q ok=%v", v, ok)
	}
}
