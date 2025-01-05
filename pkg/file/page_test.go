package file

import "testing"

const pageSize = 1024

func TestPage(t *testing.T) {
	t.Parallel()

	t.Run("Set and get int", func(t *testing.T) {
		t.Parallel()
		p := NewPage(pageSize)
		offset := 0
		value := 12345
		p.SetInt(offset, value)

		if p.GetInt(offset) != value {
			t.Errorf("expected %d, got %d", value, p.GetInt(offset))
		}
	})

	t.Run("Set and get bytes", func(t *testing.T) {
		t.Parallel()
		p := NewPage(pageSize)
		offset := 0
		value := []byte("hello, world")
		p.SetBytes(offset, value)

		if string(p.GetBytes(offset)) != string(value) {
			t.Errorf("expected %s, got %s", value, p.GetBytes(offset))
		}
	})

	t.Run("Set and get string", func(t *testing.T) {
		t.Parallel()
		p := NewPage(pageSize)
		offset := 0
		value := "hello, world"
		p.SetString(offset, value)

		if p.GetString(offset) != value {
			t.Errorf("expected %s, got %s", value, p.GetString(offset))
		}
	})

	t.Run("Max length", func(t *testing.T) {
		t.Parallel()
		expected := 4 + 10*4
		if got := MaxLength(10); got != expected {
			t.Errorf("expected %d, got %d", expected, got)
		}
	})

	t.Run("Contents", func(t *testing.T) {
		t.Parallel()
		p := NewPage(pageSize)
		offset := 0
		value := []byte("hello, world")

		p.SetBytes(offset, value)
		if p.Contents().Len() != pageSize {
			t.Errorf("expected %d, got %d", pageSize, p.Contents().Len())
		}
	})
}
