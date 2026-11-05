package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// MobType is a world.EntityType implementation for generic mobs.
type MobType struct {
	Name          string
	Width, Height float64
}

// Open converts a world.EntityHandle to an Ent.
func (t MobType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Ent{tx: tx, handle: handle, data: data}
}

// DecodeNBT decodes NBT data for the mob entity.
func (t MobType) DecodeNBT(m map[string]any, data *world.EntityData) {
	data.Data = PassiveBehaviourConfig{
		Gravity: 0.08,
		Drag:    0.02,
	}.New()
}

// EncodeNBT encodes NBT data for the mob entity.
func (t MobType) EncodeNBT(data *world.EntityData) map[string]any {
	return nil
}

// EncodeEntity returns the Minecraft entity ID.
func (t MobType) EncodeEntity() string {
	return t.Name
}

// BBox returns the bounding box of the mob entity.
func (t MobType) BBox(world.Entity) cube.BBox {
	return cube.Box(-t.Width/2, 0, -t.Width/2, t.Width/2, t.Height, t.Width/2)
}

var (
	// PigType is the EntityType for a pig.
	PigType = MobType{Name: "minecraft:pig", Width: 0.9, Height: 0.9}
	// CowType is the EntityType for a cow.
	CowType = MobType{Name: "minecraft:cow", Width: 0.9, Height: 1.4}
	// SheepType is the EntityType for a sheep.
	SheepType = MobType{Name: "minecraft:sheep", Width: 0.9, Height: 1.3}
	// ZombieType is the EntityType for a zombie.
	ZombieType = MobType{Name: "minecraft:zombie", Width: 0.6, Height: 1.9}
	// SkeletonType is the EntityType for a skeleton.
	SkeletonType = MobType{Name: "minecraft:skeleton", Width: 0.6, Height: 1.9}
)
