package components

import (
	"log"

	"github.com/pkg/errors"
)

// Type of component
type Type byte

func (t Type) String() string {
	switch t {
	case drawableType:
		return "drawableType"
	case ownedType:
		return "ownedType"
	case posType:
		return "posType"
	case rectType:
		return "rectType"
	default:
		return "unknown"
	}
}

const (
	drawableType Type = iota + 1
	ownedType
	posType
	rectType
)

// Owned is is the parent of a component
type Owned struct{}

// Drawable is a component that can be drawn
type Drawable struct{}

// Pos is the position of an entity
type Pos struct {
	X, Y int
}

// Rect is just a temp object
type Rect struct {
	W, H int
}

// Map maps between entities and components
type Map struct {
	m map[string]map[Type]interface{}
}

// NewMap returns an initialized, empty Map
func NewMap() *Map {
	m := make(map[string]map[Type]interface{})
	cm := Map{m}
	return &cm
}

func (cm *Map) add(e string, cs ...interface{}) error {

	var entry map[Type]interface{}
	// Create a entry if entity doesn't have one
	var ok bool
	if entry, ok = cm.m[e]; !ok {
		entry = make(map[Type]interface{})
	}

	failIfAleadyExist := func(t Type) {
		if _, ok := entry[t]; ok {
			log.Fatalf("Entity already has ponent of this type %d:", t)
		}
	}

	for _, c := range cs {

		var typ Type
		switch v := c.(type) {
		case Pos:
			failIfAleadyExist(typ)
			typ = posType
			entry[typ] = &v
		case Rect:
			failIfAleadyExist(typ)
			typ = rectType
			entry[typ] = &v
		case Drawable:
			failIfAleadyExist(typ)
			typ = drawableType
			entry[typ] = &v
		case Owned:
			failIfAleadyExist(typ)
			typ = ownedType
			entry[typ] = &v
		default:
			return errors.Errorf("Unknown type %v", c)
		}
	}
	cm.m[e] = entry
	return nil
}

func (cm *Map) remove(e string, typ Type) {
	delete(cm.m[e], typ)
}

func (cm *Map) removeAll(e string) {
	for key := range cm.m[e] {
		delete(cm.m[e], key)
	}
}

func (cm *Map) get(e string, typ Type) (interface{}, error) {
	if _, ok := cm.m[e]; !ok {
		return nil, errors.New("invalid ID")
	}

	return cm.m[e][typ], nil
}

func (cm *Map) hasComponents(e string, types ...Type) bool {
	if _, ok := cm.m[e]; !ok {
		return false
	}

	for _, t := range types {
		if _, ok := cm.m[e][t]; !ok {
			return false
		}
	}

	return true
}
