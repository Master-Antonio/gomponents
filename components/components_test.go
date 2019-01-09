package components

import (
	"log"
	"testing"
)

func Test_BasicComponents(t *testing.T) {

	e := "123abc"
	cm := NewMap()

	if err := cm.add(e, Pos{1, 2}); err != nil {
		t.Fatal(err)
	}

	if err := cm.add(e, Rect{3, 4}); err != nil {
		t.Fatal(err)
	}

	tcs := []struct {
		msg        string
		components []Type
		expected   bool
	}{
		{
			"Rect only",
			[]Type{rectType},
			true,
		},
		{
			"Rect and Position",
			[]Type{rectType, posType},
			true,
		},
		{
			"Drawable only",
			[]Type{Drawable},
			false,
		},
	}

	for _, tc := range tcs {
		if cm.hasComponents(e, tc.components...) != tc.expected {
			t.Fatalf("%s: expected = %t", tc.msg, tc.expected)
		}
	}

	// Updating position
	v, err := cm.get(e, posType)
	if err != nil {
		t.Fatal(err)
	}
	p := v.(*Pos)
	p.X++
	p.X++

	w, err := cm.get(e, posType)
	if err != nil {
		t.Fatal(err)
	}
	p2 := w.(*Pos)

	expected := Pos{X: 3, Y: 2}
	if *p2 != expected {
		log.Fatalf("expected %v, got %v\n", *p2, expected)
	}
}
