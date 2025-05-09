package ecs

import (
	"fmt"
	"slices"

	"github.com/lucdrenth/murph_engine/src/utils"
)

type Query interface {
	Exec(world *World) error

	// Prepare extracts the query options and puts it in CombinedQueryOptions. This should be called
	// once, after which the query is ready to be used (e.g. Exec can be called).
	//
	// This step is necessary because the way that queries are created is optimized for users, and not
	// for the computer. This method closes that gap.
	Prepare(*World) error

	// Validate checks if there are any unexpected or unoptimized things. It returns an error if there
	// is something that can be optimized. The error should be treated as a warning, and does not mean
	// that the query can not be executed.
	//
	// Prepare must be called without returning any errors before calling this method.
	Validate() error

	// Clear the query results that got filled when last running Exec.
	ClearResults()

	// IsLazy returns wether this query should be treated as lazy or not. Being lazy means that it should
	// not get executed automatically, and should be done by the user.
	IsLazy() bool

	getOptions() *combinedQueryOptions
}

type queryOptions struct {
	options    combinedQueryOptions
	components []ComponentId
}

func (o *queryOptions) getOptions() *combinedQueryOptions {
	return &o.options
}

func (o *queryOptions) IsLazy() bool {
	return o.options.isLazy
}

func (o *queryOptions) Validate() error {
	return o.options.validateOptions(o.components)
}

// Query1 queries 1 component.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query1[ComponentA IComponent, _ QueryOption] struct {
	queryOptions
	results Query1Result[ComponentA]
}

// Query2 queries 2 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query2[ComponentA, ComponentB IComponent, _ QueryOption] struct {
	queryOptions
	results Query2Result[ComponentA, ComponentB]
}

// Query3 queries 3 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query3[ComponentA, ComponentB, ComponentC IComponent, _ QueryOption] struct {
	queryOptions
	results Query3Result[ComponentA, ComponentB, ComponentC]
}

// Query4 queries 4 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query4[ComponentA, ComponentB, ComponentC, ComponentD IComponent, _ QueryOption] struct {
	queryOptions
	results Query4Result[ComponentA, ComponentB, ComponentC, ComponentD]
}

func (q *Query1[ComponentA, QueryOptions]) Exec(world *World) error {
	q.ClearResults()

	for entityId, entityData := range world.entities {
		if q.options.isFilteredOut(entityData) {
			continue
		}

		a, match, err := getQueryComponent[ComponentA](world, entityData, &q.options)
		if err != nil {
			return err
		}
		if !match {
			continue
		}

		q.results.componentsA = append(q.results.componentsA, a)
		q.results.entityIds = append(q.results.entityIds, entityId)
	}

	return nil
}
func (q *Query2[ComponentA, ComponentB, QueryOptions]) Exec(world *World) error {
	q.ClearResults()

	for entityId, entityData := range world.entities {
		if q.options.isFilteredOut(entityData) {
			continue
		}

		a, match, err := getQueryComponent[ComponentA](world, entityData, &q.options)
		if err != nil {
			return err
		}
		if !match {
			continue
		}

		b, match, err := getQueryComponent[ComponentB](world, entityData, &q.options)
		if err != nil {
			return err
		}
		if !match {
			continue
		}

		q.results.componentsA = append(q.results.componentsA, a)
		q.results.componentsB = append(q.results.componentsB, b)
		q.results.entityIds = append(q.results.entityIds, entityId)
	}

	return nil
}
func (q *Query3[ComponentA, ComponentB, ComponentC, QueryOptions]) Exec(world *World) error {
	q.ClearResults()

	for entityId, entityData := range world.entities {
		if q.options.isFilteredOut(entityData) {
			continue
		}

		a, match, err := getQueryComponent[ComponentA](world, entityData, &q.options)
		if err != nil {
			return err
		}
		if !match {
			continue
		}

		b, match, err := getQueryComponent[ComponentB](world, entityData, &q.options)
		if err != nil {
			return err
		}
		if !match {
			continue
		}

		c, match, err := getQueryComponent[ComponentC](world, entityData, &q.options)
		if err != nil {
			return err
		}
		if !match {
			continue
		}

		q.results.componentsA = append(q.results.componentsA, a)
		q.results.componentsB = append(q.results.componentsB, b)
		q.results.componentsC = append(q.results.componentsC, c)
		q.results.entityIds = append(q.results.entityIds, entityId)
	}

	return nil
}
func (q *Query4[ComponentA, ComponentB, ComponentC, ComponentD, QueryOptions]) Exec(world *World) error {
	q.ClearResults()

	for entityId, entityData := range world.entities {
		if q.options.isFilteredOut(entityData) {
			continue
		}

		a, match, err := getQueryComponent[ComponentA](world, entityData, &q.options)
		if err != nil {
			return err
		}
		if !match {
			continue
		}

		b, match, err := getQueryComponent[ComponentB](world, entityData, &q.options)
		if err != nil {
			return err
		}
		if !match {
			continue
		}

		c, match, err := getQueryComponent[ComponentC](world, entityData, &q.options)
		if err != nil {
			return err
		}
		if !match {
			continue
		}

		d, match, err := getQueryComponent[ComponentD](world, entityData, &q.options)
		if err != nil {
			return err
		}
		if !match {
			continue
		}

		q.results.componentsA = append(q.results.componentsA, a)
		q.results.componentsB = append(q.results.componentsB, b)
		q.results.componentsC = append(q.results.componentsC, c)
		q.results.componentsD = append(q.results.componentsD, d)
		q.results.entityIds = append(q.results.entityIds, entityId)
	}

	return nil
}

func (q *Query1[A, QueryOptions]) Prepare(world *World) (err error) {
	q.options, err = toCombinedQueryOptions[QueryOptions](world)
	q.components = []ComponentId{
		ComponentIdFor[A](world),
	}
	return err
}
func (q *Query2[A, B, QueryOptions]) Prepare(world *World) (err error) {
	q.options, err = toCombinedQueryOptions[QueryOptions](world)
	q.components = []ComponentId{
		ComponentIdFor[A](world),
		ComponentIdFor[B](world),
	}
	return err
}
func (q *Query3[A, B, C, QueryOptions]) Prepare(world *World) (err error) {
	q.options, err = toCombinedQueryOptions[QueryOptions](world)
	q.components = []ComponentId{
		ComponentIdFor[A](world),
		ComponentIdFor[B](world),
		ComponentIdFor[C](world),
	}
	return err
}
func (q *Query4[A, B, C, D, QueryOptions]) Prepare(world *World) (err error) {
	q.options, err = toCombinedQueryOptions[QueryOptions](world)
	q.components = []ComponentId{
		ComponentIdFor[A](world),
		ComponentIdFor[B](world),
		ComponentIdFor[C](world),
		ComponentIdFor[D](world),
	}
	return err
}

func (q *Query1[ComponentA, QueryOptions]) Result() *Query1Result[ComponentA] {
	return &q.results
}
func (q *Query2[ComponentA, ComponentB, QueryOptions]) Result() *Query2Result[ComponentA, ComponentB] {
	return &q.results
}
func (q *Query3[ComponentA, ComponentB, ComponentC, QueryOptions]) Result() *Query3Result[ComponentA, ComponentB, ComponentC] {
	return &q.results
}
func (q *Query4[ComponentA, ComponentB, ComponentC, ComponentD, QueryOptions]) Result() *Query4Result[ComponentA, ComponentB, ComponentC, ComponentD] {
	return &q.results
}

func (q *Query1[ComponentA, QueryOptions]) ClearResults() {
	q.results.Clear()
}
func (q *Query2[ComponentA, ComponentB, QueryOptions]) ClearResults() {
	q.results.Clear()
}
func (q *Query3[ComponentA, ComponentB, ComponentC, QueryOptions]) ClearResults() {
	q.results.Clear()
}
func (q *Query4[ComponentA, ComponentB, ComponentC, ComponentD, QueryOptions]) ClearResults() {
	q.results.Clear()
}

// getQueryComponent returns a pointer to T if the component is found on the entity.
//
// match is true when the entity has the component or if the component is marked marked as optional.
// When match is true, the entity should be present in the query results.
func getQueryComponent[T IComponent](world *World, entityData *EntityData, queryOptions *combinedQueryOptions) (result *T, match bool, err error) {
	componentId := ComponentIdFor[T](world)

	componentRegistryIndex, entityHasComponent := entityData.components[componentId]
	if !entityHasComponent {
		return nil, slices.Contains(queryOptions.OptionalComponents, componentId), nil
	}

	result, err = getComponentFromComponentRegistry[T](world.components[componentId], componentRegistryIndex)
	if err != nil {
		return nil, false, fmt.Errorf("getQueryComponent encountered unexpected error: %v", err)
	}

	if result != nil && (queryOptions.ReadOnlyComponents.IsAllReadOnly || slices.Contains(queryOptions.ReadOnlyComponents.ComponentIds, componentId)) {
		result = utils.ClonePointerValue(result)
	}

	return result, true, nil
}
