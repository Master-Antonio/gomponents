package direction

import (
	"strings"

	"github.com/peterhellberg/gfx"
)

type D byte

func (d D) String() string {
	var s string
	if d&Up == Up {
		s += "Up "
	}
	if d&Down == Down {
		s += "Down "
	}
	if d&Left == Left {
		s += "Left "
	}
	if d&Right == Right {
		s += "Right "
	}

	return s
}

var (
	Up    D = D(1 << 0)
	Down  D = D(1 << 1)
	Left  D = D(1 << 2)
	Right D = D(1 << 3)
	UDLR  D = Up | Down | Left | Right
	All   D = UDLR
)

func FromVec(v gfx.Vec) D {
	var dir D
	if v.X > 0 {
		dir |= Right
	}

	if v.X < 0 {
		dir |= Left
	}

	if v.Y > 0 {
		dir |= Up
	}

	if v.Y < 0 {
		dir |= Down
	}
	return dir
}

func FromString(s string) D {
	if s == "" {
		return D(Up | Down | Left | Right)
	}

	s = strings.ToLower(s)

	var dir D
	if strings.Contains(s, "u") {
		dir |= Up
	}
	if strings.Contains(s, "d") {
		dir |= Down
	}
	if strings.Contains(s, "l") {
		dir |= Left
	}
	if strings.Contains(s, "r") {
		dir |= Right
	}
	return dir
}
