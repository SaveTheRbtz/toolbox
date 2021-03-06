package mo_string

import "testing"

func TestSelectString(t *testing.T) {
	s1 := NewSelect()
	s1.SetOptions("orange", "apple", "orange", "grape")
	if !s1.IsValid() {
		t.Error(s1.IsValid())
	}
	if s1.Value() != "orange" {
		t.Error(s1.Value())
	}
	s1.SetSelect("pine")
	if s1.IsValid() {
		t.Error(s1.IsValid())
	}
	if len(s1.Options()) != 3 {
		t.Error(s1.Options())
	}
}
