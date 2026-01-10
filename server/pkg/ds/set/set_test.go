package set

import "testing"

func TestSet(t *testing.T) {
	s := NewSet[string]()

	if !s.IsEmpty() {
		t.Error("new set must be empty")
	}

	ss := NewSet("1", "2")

	if ss.Size() != 2 {
		t.Errorf("Expected error %v, got %v", s.Size(), 2)
	}
}

func TestAdd(t *testing.T) {
	s := NewSet[int]()

	s.Add(1)
	if !s.Contain(1) {
		t.Error("set should contain 1 after adding")
	}
	if s.Size() != 1 {
		t.Errorf("expected size 1, got %d", s.Size())
	}

	s.Add(2)
	s.Add(3)
	if s.Size() != 3 {
		t.Errorf("expected size 3, got %d", s.Size())
	}

	s.Add(1)
	if s.Size() != 3 {
		t.Error("adding duplicate should not increase size")
	}
}

func TestRemove(t *testing.T) {
	s := NewSet(1, 2, 3)

	s.Remove(1)
	if s.Contain(1) {
		t.Error("set should not contain 1 after removing")
	}
	if s.Size() != 2 {
		t.Errorf("expected size 2, got %d", s.Size())
	}

	s.Remove(2)
	s.Remove(3)
	if !s.IsEmpty() {
		t.Error("set should be empty after removing all items")
	}

	s.Remove(999)
	if s.Size() != 0 {
		t.Error("removing non-existent item should not affect size")
	}
}

func TestClear(t *testing.T) {
	s := NewSet(1, 2, 3)

	s.Clear()

	if s.Size() != 0 {
		t.Error("set should be empty after clear")
	}
}

func TestItems(t *testing.T) {
	s := NewSet(1, 2, 3)

	items := s.Items()

	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d", len(items))
	}

	itemMap := make(map[int]bool, s.Size())
	for _, item := range items {
		itemMap[item] = true
	}

	if !itemMap[1] || !itemMap[2] || !itemMap[3] {
		t.Error("items should contain 1, 2, and 3")
	}
}
