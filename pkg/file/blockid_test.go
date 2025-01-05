package file

import "testing"

func TestBlockID(t *testing.T) {
	t.Parallel()

	b1 := NewBlockID("file1", 1)
	b2 := NewBlockID("file1", 2)
	b3 := NewBlockID("file2", 1)

	if b1.Equals(b2) {
		t.Errorf("b1 should not equal b2")
	}

	if b1.Equals(b3) {
		t.Errorf("b1 should not equal b3")
	}

	if !b1.Equals(b1) {
		t.Errorf("b1 should equal itself")
	}

	if b1.HashCode() == b2.HashCode() {
		t.Errorf("b1 and b2 should have different hash codes")
	}

	if b1.HashCode() == b3.HashCode() {
		t.Errorf("b1 and b3 should have different hash codes")
	}

}
