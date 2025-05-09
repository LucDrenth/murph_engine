package ecs

import (
	"errors"
	"reflect"
)

// World contains all of the entities and their components.
type World struct {
	entityIdCounter          uint
	entities                 map[EntityId]*EntityData
	components               map[ComponentId]*componentRegistry
	initialComponentCapacity initialComponentCapacityStrategy
	componentIdRegistry      componentIdRegistry
}

// DefaultWorld returns a World with default configs.
func DefaultWorld() World {
	world, err := NewWorld(DefaultWorldConfigs())
	if err != nil {
		// Creating a world with default configs should never result in an error.
		// This is confirmed by a unit test, so we can 'safely' panic just in case
		// it happens.
		panic(err)
	}

	return world
}

// NewWorld returns a world that can contain entities and components.
func NewWorld(configs WorldConfigs) (World, error) {
	if configs.ComponentCapacityStrategy == nil {
		return World{}, errors.New("component capacity strategy can not be nil")
	}

	return World{
		entities:                 map[EntityId]*EntityData{},
		components:               map[ComponentId]*componentRegistry{},
		initialComponentCapacity: configs.ComponentCapacityStrategy,
		componentIdRegistry: componentIdRegistry{
			components: map[reflect.Type]uint{},
		},
	}, nil
}

func (world *World) CountEntities() int {
	return len(world.entities)
}

func (world *World) CountComponents() int {
	result := 0

	for _, entry := range world.entities {
		result += len(entry.components)
	}

	return result
}

func (world *World) createEntity() EntityId {
	world.entityIdCounter++
	entity := EntityId(world.entityIdCounter)
	world.entities[entity] = &EntityData{components: map[ComponentId]uint{}}
	return entity
}

// getComponentRegistry creates a new component registry if it doesn't exist yet.
func (world *World) getComponentRegistry(ComponentId ComponentId) (*componentRegistry, error) {
	componentRegistry, ok := world.components[ComponentId]
	if !ok {
		newComponentRegistry, err := createComponentRegistry(
			world.initialComponentCapacity.GetDefaultComponentCapacity(ComponentId),
			ComponentId,
		)

		if err != nil {
			return nil, err
		}

		world.components[ComponentId] = &newComponentRegistry
		return &newComponentRegistry, nil
	}

	return componentRegistry, nil
}
