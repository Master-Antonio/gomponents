package components

import (
	"fmt"
	"log"
	"os"
	"testing"
	"text/tabwriter"
)

type testCase struct {
	msg      string
	types    []Type
	expected bool
}

func Test_Add(t *testing.T) {
	e := "123abc"
	cm := NewMap()

	tcs := []interface{}{Pos{1, 2}, Rect{3, 4}}

	for _, tc := range tcs {
		if err := cm.add(e, tc); err != nil {
			t.Fatal(err)
		}
	}
}

func Test_Set(t *testing.T) {

	e := "123abc"
	cm := NewMap()

	if err := cm.add(e, Pos{1, 2}); err != nil {
		t.Fatal(err)
	}

	// Updating position
	v, err := cm.get(e, posType)
	if err != nil {
		t.Fatal(err)
	}
	p := v.(*Pos)
	p.X++
	p.X++

	// Get position again
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

func Test_HasComponents(t *testing.T) {
	e := "123abc"
	cm := NewMap()

	if err := cm.add(e, Pos{1, 2}); err != nil {
		t.Fatal(err)
	}

	if err := cm.add(e, Rect{3, 4}); err != nil {
		t.Fatal(err)
	}

	tcs := []testCase{
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
			[]Type{drawableType},
			false,
		},
	}

	for _, tc := range tcs {
		if cm.hasComponents(e, tc.types...) != tc.expected {
			t.Fatalf("%s: expected = %t", tc.msg, tc.expected)
		}
	}
}

func Test_Remove(t *testing.T) {

	e := "123abc"
	cm := NewMap()

	if err := cm.add(e, Pos{1, 2}); err != nil {
		t.Fatal(err)
	}
	if err := cm.add(e, Drawable{}); err != nil {
		t.Fatal(err)
	}

	tcs := []testCase{
		{"", []Type{posType}, true},
		{"", []Type{drawableType}, true},
		{"", []Type{ownedType}, false},
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "HAS\tEXPECTED\tGOT\t\n")
	for _, tc := range tcs {
		actual := cm.hasComponents(e, tc.types...)
		fmt.Fprintf(w, "%v\t%t\t%t\t\n", tc.types, tc.expected, actual)

		if actual != tc.expected {
			w.Flush()
			t.FailNow()
		}
	}

	// Remove position
	cm.remove(e, posType)

	// Remove non-existing type
	cm.remove(e, ownedType) // Shouldn't crash

	// Remove from non-existing entity
	cm.remove("xyz987", ownedType) // Shouldn't crash

	tcs = []testCase{
		{"", []Type{posType}, false},
		{"", []Type{drawableType}, true},
		{"", []Type{ownedType}, false},
	}
	w = tabwriter.NewWriter(os.Stdout, 12, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "HAS\tEXPECTED\tGOT\t\n")
	for _, tc := range tcs {
		actual := cm.hasComponents(e, tc.types...)
		fmt.Fprintf(w, "%v\t%t\t%t\t\n", tc.types, tc.expected, actual)

		if actual != tc.expected {
			w.Flush()
			t.FailNow()
		}
	}
}

func Test_RemoveAll(t *testing.T) {

	e := "abc123"
	cm := NewMap()

	if err := cm.add(e, Pos{}); err != nil {
		t.Fatal(err)
	}
	if err := cm.add(e, Owned{}); err != nil {
		t.Fatal(err)
	}

	cm.removeAll(e)
	if cm.hasComponents(e, posType) || cm.hasComponents(e, ownedType) || cm.hasComponents(e, drawableType) {
		t.Fatalf("Expected no components, got %s=%t,%s=%t,%s=%t",
			posType, cm.hasComponents(e, posType),
			ownedType, cm.hasComponents(e, ownedType),
			drawableType, cm.hasComponents(e, drawableType))
	}
}
