package components

import (
	"fmt"
	"log"
	"os"
	"testing"
	"text/tabwriter"

	"github.com/peterhellberg/gfx"
)

type testCase struct {
	msg      string
	types    []Type
	expected bool
}

func Test_Add(t *testing.T) {
	e := "123abc"
	cm := NewMap()

	tcs := []interface{}{Pos{gfx.V(1, 2)}}

	for _, tc := range tcs {
		if err := cm.Add(e, tc); err != nil {
			t.Fatal(err)
		}
	}
}

func Test_Set(t *testing.T) {

	e := "123abc"
	cm := NewMap()

	if err := cm.Add(e, Pos{gfx.V(1, 2)}); err != nil {
		t.Fatal(err)
	}

	// Updating position
	v, err := cm.Get(e, PosType)
	if err != nil {
		t.Fatal(err)
	}
	p := v.(*Pos)
	p.X++
	p.X++

	// Get position again
	w, err := cm.Get(e, PosType)
	if err != nil {
		t.Fatal(err)
	}
	p2 := w.(*Pos)

	expected := Pos{gfx.V(3, 2)}
	if *p2 != expected {
		log.Fatalf("expected %v, got %v\n", *p2, expected)
	}
}

func Test_HasComponents(t *testing.T) {
	e := "123abc"
	cm := NewMap()

	if err := cm.Add(e, Pos{gfx.V(1, 2)}); err != nil {
		t.Fatal(err)
	}

	if err := cm.Add(e, Owned{}); err != nil {
		t.Fatal(err)
	}

	tcs := []testCase{
		{
			"Owned only",
			[]Type{OwnedType},
			true,
		},
		{
			"Owned and Position",
			[]Type{OwnedType, PosType},
			true,
		},
		{
			"Drawable only",
			[]Type{DrawableType},
			false,
		},
	}

	for _, tc := range tcs {
		if cm.HasComponents(e, tc.types...) != tc.expected {
			t.Fatalf("%s: expected = %t", tc.msg, tc.expected)
		}
	}
}

func Test_Remove(t *testing.T) {

	e := "123abc"
	cm := NewMap()

	if err := cm.Add(e, Pos{gfx.V(1, 2)}); err != nil {
		t.Fatal(err)
	}
	if err := cm.Add(e, Drawable{}); err != nil {
		t.Fatal(err)
	}

	tcs := []testCase{
		{"", []Type{PosType}, true},
		{"", []Type{DrawableType}, true},
		{"", []Type{OwnedType}, false},
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "HAS\tEXPECTED\tGOT\t\n")
	for _, tc := range tcs {
		actual := cm.HasComponents(e, tc.types...)
		fmt.Fprintf(w, "%v\t%t\t%t\t\n", tc.types, tc.expected, actual)

		if actual != tc.expected {
			w.Flush()
			t.FailNow()
		}
	}

	// Remove position
	cm.Remove(e, PosType)

	// Remove non-existing type
	cm.Remove(e, OwnedType) // Shouldn't crash

	// Remove from non-existing entity
	cm.Remove("xyz987", OwnedType) // Shouldn't crash

	tcs = []testCase{
		{"", []Type{PosType}, false},
		{"", []Type{DrawableType}, true},
		{"", []Type{OwnedType}, false},
	}
	w = tabwriter.NewWriter(os.Stdout, 12, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "HAS\tEXPECTED\tGOT\t\n")
	for _, tc := range tcs {
		actual := cm.HasComponents(e, tc.types...)
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

	if err := cm.Add(e, Pos{}); err != nil {
		t.Fatal(err)
	}
	if err := cm.Add(e, Owned{}); err != nil {
		t.Fatal(err)
	}

	cm.RemoveAll(e)
	if cm.HasComponents(e, PosType) || cm.HasComponents(e, OwnedType) || cm.HasComponents(e, DrawableType) {
		t.Fatalf("Expected no components, got %s=%t,%s=%t,%s=%t",
			PosType, cm.HasComponents(e, PosType),
			OwnedType, cm.HasComponents(e, OwnedType),
			DrawableType, cm.HasComponents(e, DrawableType))
	}
}
