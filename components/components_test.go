package components

import (
	"fmt"
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
			[]Type{RectType},
			true,
		},
		{
			"Rect and Position",
			[]Type{RectType, PositionType},
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

	v, err := cm.get(e, RectType)
	if err != nil {
		t.Fatal(err)
	}
	r := (&v).(Rect)
	fmt.Println(r)

	// fmt.Println("has Rect, Pos=", cm.hasComponents(e, RectType))
	// fmt.Println("has Rect, Pos=", cm.hasComponents(e, RectType, PositionType))

	// componentIndex().addComponent<ComponentTypeT>(m_id, std::forward<Args>(args)...);

}
