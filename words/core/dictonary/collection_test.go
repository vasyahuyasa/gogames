package dictonary

import (
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	c := New()

	if c.HasKey("test") || c.HasKey("") || c.HasKey("a") {
		t.Error("HasKey вернула true несуществующему ключу")
	}
}

func TestSet(t *testing.T) {
	c := New()

	test := []string{"test", "", "a", "gopher"}

	for _, key := range test {
		c.SetKey(key)
	}

	for _, key := range test {
		if !c.HasKey(key) {
			t.Error("HasKey вернула false существующему ключу", key)
		}
	}
}

func TestDelete(t *testing.T) {
	c := New()

	test := []string{"test", "", "a", "gopher"}
	for _, key := range test {
		c.SetKey(key)
	}

	for _, key := range test {
		c.DeleteKey(key)
	}

	for _, key := range test {
		if c.HasKey(key) {
			t.Error("HasKey вернула true несуществующему ключу", key)
		}
	}
}

func TestReset(t *testing.T) {
	c := New()

	test := []string{"test", "", "a", "gopher"}
	for _, key := range test {
		c.SetKey(key)
	}

	c.Rest()

	for _, key := range test {
		if c.HasKey(key) {
			t.Error("HasKey вернула true несуществующему ключу", key)
		}
	}
}

func TestReader(t *testing.T) {
	test := []string{"test", "test2", "hi", "gopher", "города"}
	r := strings.NewReader(strings.Join(test, "\n"))
	c := New()
	err := c.FromReader(r)
	if err != nil {
		t.Errorf("Чтение из reader: %v", err)
	}

	for _, key := range test {
		if !c.HasKey(key) {
			t.Error("HasKey вернула false существующему ключу", key)
		}
	}
}
