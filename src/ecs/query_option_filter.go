package ecs

import (
	"errors"

	"github.com/lucdrenth/murph_engine/src/utils"
)

type filterType = int

const (
	filterTypeWith filterType = iota
	filterTypeWithout
	filterTypeAnd
	filterTypeOr
	filterTypeNone
)

type QueryParamFilter interface {
	getComponents(world *World) []ComponentId
	getFilterType() filterType
	getNestedFilters() (a QueryParamFilter, b QueryParamFilter, err error)
}
type NoFilter struct{}
type With[A IComponent] struct{}
type Without[A IComponent] struct{}
type And[A, B QueryParamFilter] struct{}
type Or[A, B QueryParamFilter] struct{}

func (filter NoFilter) getComponents(world *World) []ComponentId {
	return []ComponentId{}
}

func (filter And[A, B]) getComponents(world *World) []ComponentId {
	return []ComponentId{}
}

func (filter Or[A, B]) getComponents(world *World) []ComponentId {
	return []ComponentId{}
}

func (filter With[A]) getComponents(world *World) []ComponentId {
	return []ComponentId{ComponentIdFor[A](world)}
}

func (filter Without[A]) getComponents(world *World) []ComponentId {
	return []ComponentId{ComponentIdFor[A](world)}
}

func (filter NoFilter) getFilterType() filterType {
	return filterTypeNone
}

func (filter With[A]) getFilterType() filterType {
	return filterTypeWith
}

func (filter Without[A]) getFilterType() filterType {
	return filterTypeWithout
}

func (filter And[A, B]) getFilterType() filterType {
	return filterTypeAnd
}

func (filter Or[A, B]) getFilterType() filterType {
	return filterTypeOr
}

func (filter NoFilter) getNestedFilters() (a QueryParamFilter, b QueryParamFilter, err error) {
	return nil, nil, errors.New("nested filters not supported for this type")
}

func (filter With[A]) getNestedFilters() (a QueryParamFilter, b QueryParamFilter, err error) {
	return nil, nil, errors.New("nested filters not supported for this type")

}

func (filter Without[A]) getNestedFilters() (a QueryParamFilter, b QueryParamFilter, err error) {
	return nil, nil, errors.New("nested filters not supported for this type")
}

func (filter And[A, B]) getNestedFilters() (a QueryParamFilter, b QueryParamFilter, err error) {
	a, err = utils.ToConcrete[A]()
	if err != nil {
		return nil, nil, err
	}

	b, err = utils.ToConcrete[B]()
	if err != nil {
		return nil, nil, err
	}

	return a, b, nil
}

func (filter Or[A, B]) getNestedFilters() (a QueryParamFilter, b QueryParamFilter, err error) {
	a, err = utils.ToConcrete[A]()
	if err != nil {
		return nil, nil, err
	}

	b, err = utils.ToConcrete[B]()
	if err != nil {
		return nil, nil, err
	}

	return a, b, nil
}

func (filter With[A]) getCombinedQueryOptions(world *World) (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[With[A], NoOptional, NoReadOnly, NotLazy]](world)
}
func (filter Without[A]) getCombinedQueryOptions(world *World) (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[Without[A], NoOptional, NoReadOnly, NotLazy]](world)
}
func (filter And[A, B]) getCombinedQueryOptions(world *World) (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[And[A, B], NoOptional, NoReadOnly, NotLazy]](world)
}
func (filter Or[A, B]) getCombinedQueryOptions(world *World) (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[Or[A, B], NoOptional, NoReadOnly, NotLazy]](world)
}

type QueryFilter interface {
	// Validate that entityData satisfies the filter
	Validate(*EntityData) bool
}
type queryFilterAnd struct {
	a QueryFilter
	b QueryFilter
}
type queryFilterOr struct {
	a QueryFilter
	b QueryFilter
}
type queryFilterWith struct {
	c []ComponentId
}
type queryFilterWithout struct {
	c []ComponentId
}

func (filter queryFilterAnd) Validate(e *EntityData) bool {
	return filter.a.Validate(e) && filter.b.Validate(e)
}

func (filter queryFilterOr) Validate(e *EntityData) bool {
	return filter.a.Validate(e) || filter.b.Validate(e)
}

func (filter queryFilterWith) Validate(e *EntityData) bool {
	for _, c := range filter.c {
		if !e.hasComponent(c) {
			return false
		}
	}

	return true
}

func (filter queryFilterWithout) Validate(e *EntityData) bool {
	for _, c := range filter.c {
		if e.hasComponent(c) {
			return false
		}
	}

	return true
}
