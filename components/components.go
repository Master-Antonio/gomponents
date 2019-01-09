package components

import (
	"fmt"

	"github.com/pkg/errors"
)

// Type of component
type Type byte

const (
	// Drawable objects can be drawn to an image
	Drawable Type = iota + 1
	// Owned by something
	Owned

	// PosType of something
	posType

	rectType
)

type Pos struct {
	X, Y int
}

type Rect struct {
	W, H int
}

// TypeID is the ID in the corresponding map
type TypeID int32

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

	for _, c := range cs {

		var typ Type
		switch v := c.(type) {
		case Pos:
			fmt.Println("Pos")
			typ = posType
			entry[typ] = &v
		case Rect:
			fmt.Println("Rect")
			typ = rectType
			entry[typ] = &v
		default:
			return errors.Errorf("Unknown type %v", c)
		}

		// Check that entity doesn't have component
		// if _, ok := entry[typ]; ok {
		// 	log.Fatalf("Entity already has ponent of this type %d:", typ)
		// 	// return errors.Errorf("Entity acomlready has component of this type %d:", typ)
		// }

		fmt.Println("pointer", &c, c)
	}
	cm.m[e] = entry
	return nil
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

func (cm *Map) get(e string, typ Type) (interface{}, error) {
	if _, ok := cm.m[e]; !ok {
		return nil, errors.New("invalid ID")
	}

	return cm.m[e][typ], nil
}
