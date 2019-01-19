package components

import (
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
	ase "github.com/kyeett/GoAseprite"
	"github.com/kyeett/gomponents/direction"
	"github.com/peterhellberg/gfx"
	"github.com/pkg/errors"
)

// Type of component
type Type byte

func (t Type) String() string {
	switch t {
	case DrawableType:
		return "DrawableType"
	case OwnedType:
		return "OwnedType"
	case PosType:
		return "PosType"
	case VelocityType:
		return "VelocityType"
	case SpriteType:
		return "SpriteType"
	case HitboxType:
		return "HitboxType"
	case AnimatedType:
		return "AnimatedType"
	default:
		return "unknown"
	}
}

// ComponentTypes that are used for the Get/Set/HasComponents calls
const (
	DrawableType Type = iota + 1
	HitboxType
	OwnedType
	PosType
	VelocityType
	SpriteType
	AnimatedType
	DirectionType
	TagsType
	CounterType
	HazardType
	BouncyType
	KillableType
	ScenarioType
	RotatedType
	TeleportingType
	TriggerType
	FollowingType
)

// Owned is is the parent of a component
type Owned struct{}

// Drawable is a component that can be drawn
type Drawable struct {
	*ebiten.Image
}

// Pos is the position of an entity
type Pos struct {
	gfx.Vec
}

// Following follows entity with ID at an offset
type Following struct {
	ID     string
	Offset gfx.Vec
}

// Velocity of an entity
type Velocity struct {
	gfx.Vec
}

// Hitbox is the rectangle used for collisions
type Hitbox struct {
	gfx.Rect
	Properties map[string]bool
	Target     string
}

// NewHitbox returns a new hitbox
func NewHitbox(rect gfx.Rect) Hitbox {
	return Hitbox{
		Rect:       rect,
		Properties: make(map[string]bool),
		Target:     "",
	}
}

// Trigger is a an area that triggers a scenario
type Trigger struct {
	gfx.Rect
	Scenario  string
	Direction direction.D
}

func (t *Trigger) String() string {
	return fmt.Sprintf("Trigger: %s at %v from %s", t.Scenario, t.Rect, t.Direction)
}

// Animated contains Aseprite sprite animation information
type Animated struct {
	Ase ase.File
}

// Direction of an entity
type Direction struct {
	D float64
}

// Hazard marks an entity as dangerous
type Hazard struct{}

// Bouncy marks an entity as bouncy
type Bouncy struct{}

// Killable marks an entity as killable
type Killable struct{}

// Teleporting teleports entity on collision
type Teleporting struct {
	Name, Target string
	Pos          gfx.Vec
}

// Rotated denotes the rotation of an entity
type Rotated struct {
	Angle float64
}

// Rotate angle further. Wraps at 2*Pi
func (r *Rotated) Rotate(v float64) {
	r.Angle = math.Mod(r.Angle+v, 2*math.Pi)
}

// Counter is a map of string keys and values
type Counter map[string]int

// Scenario is a function that gets called every turn
type Scenario struct {
	F func() bool
}

// Tags is a generic strings that can be attached to an entity
type Tags []string

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

// Add component typ cs to entity e
func (cm *Map) Add(e string, cs ...interface{}) error {

	var entry map[Type]interface{}
	// Create a entry if entity doesn't have one
	var ok bool
	if entry, ok = cm.m[e]; !ok {
		entry = make(map[Type]interface{})
	}

	failIfAlreadyExist := func(t Type) {
		if _, ok := entry[t]; ok {
			log.Fatalf("Entity already has ponent of this type %d:", t)
		}
	}

	for _, c := range cs {

		var typ Type
		switch v := c.(type) {
		case Pos:
			failIfAlreadyExist(typ)
			typ = PosType
			entry[typ] = &v
		case Drawable:
			failIfAlreadyExist(typ)
			typ = DrawableType
			entry[typ] = &v
		case Owned:
			failIfAlreadyExist(typ)
			typ = OwnedType
			entry[typ] = &v
		case Tags:
			failIfAlreadyExist(typ)
			typ = TagsType
			entry[typ] = &v
		case Hitbox:
			failIfAlreadyExist(typ)
			typ = HitboxType
			entry[typ] = &v
		case Velocity:
			failIfAlreadyExist(typ)
			typ = VelocityType
			entry[typ] = &v
		case Counter:
			failIfAlreadyExist(typ)
			typ = CounterType
			entry[typ] = &v
		case Animated:
			failIfAlreadyExist(typ)
			typ = AnimatedType
			entry[typ] = &v
		case Direction:
			failIfAlreadyExist(typ)
			typ = DirectionType
			entry[typ] = &v
		case Hazard:
			failIfAlreadyExist(typ)
			typ = HazardType
			entry[typ] = &v
		case Bouncy:
			failIfAlreadyExist(typ)
			typ = BouncyType
			entry[typ] = &v
		case Killable:
			failIfAlreadyExist(typ)
			typ = KillableType
			entry[typ] = &v
		case Scenario:
			failIfAlreadyExist(typ)
			typ = ScenarioType
			entry[typ] = &v
		case Rotated:
			failIfAlreadyExist(typ)
			typ = RotatedType
			entry[typ] = &v
		case Teleporting:
			failIfAlreadyExist(typ)
			typ = TeleportingType
			entry[typ] = &v
		case Trigger:
			failIfAlreadyExist(typ)
			typ = TriggerType
			entry[typ] = &v
		case Following:
			failIfAlreadyExist(typ)
			typ = FollowingType
			entry[typ] = &v
		default:
			return errors.Errorf("Unknown type %v", c)
		}
	}
	cm.m[e] = entry
	return nil
}

// Remove component of type typ for entity e
func (cm *Map) Remove(e string, typ Type) {
	delete(cm.m[e], typ)
}

// RemoveAll components for entity e
func (cm *Map) RemoveAll(e string) {
	for key := range cm.m[e] {
		delete(cm.m[e], key)
	}
}

// Get component of type typ for entity e
func (cm *Map) Get(e string, typ Type) (interface{}, error) {
	if _, ok := cm.m[e]; !ok {
		return nil, errors.New("invalid ID")
	}
	return cm.m[e][typ], nil
}

// GetUnsafe gets an object without checking the error. Use HasComponents before, to check multiple types
func (cm *Map) GetUnsafe(e string, typ Type) interface{} {
	return cm.m[e][typ]
}

// HasComponents returns true of entity e has all the componens of type types
func (cm *Map) HasComponents(e string, types ...Type) bool {
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
