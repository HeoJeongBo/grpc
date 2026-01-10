package item

import (
	"strings"

	itemv1 "grpc-server/proto-generated/item"
)

type FilterFunc func(*itemv1.Item) bool

type FilterBuilder struct {
	filters []FilterFunc
}

func NewFilterBuilder() *FilterBuilder {
	return &FilterBuilder{
		filters: make([]FilterFunc, 0),
	}
}

func (fb *FilterBuilder) AddFilter(f FilterFunc) *FilterBuilder {
	if f != nil {
		fb.filters = append(fb.filters, f)
	}
	return fb
}

func (fb *FilterBuilder) Build() FilterFunc {
	if len(fb.filters) == 0 {
		return func(*itemv1.Item) bool { return true }
	}

	return func(item *itemv1.Item) bool {
		for _, filter := range fb.filters {
			if !filter(item) {
				return false
			}
		}
		return true
	}
}

func NameFilter(name string) FilterFunc {
	if name == "" {
		return nil
	}

	lowerName := strings.ToLower(name)
	return func(item *itemv1.Item) bool {
		return strings.Contains(strings.ToLower(item.Name), lowerName)
	}
}

func DescriptionFilter(description string) FilterFunc {
	if description == "" {
		return nil
	}

	lowerDesc := strings.ToLower(description)
	return func(item *itemv1.Item) bool {
		return strings.Contains(strings.ToLower(item.Description), lowerDesc)
	}
}

func StatusFilter(statuses []itemv1.ItemStatus) FilterFunc {
	if len(statuses) == 0 {
		return nil
	}

	statusMap := make(map[itemv1.ItemStatus]bool, len(statuses))
	for _, status := range statuses {
		statusMap[status] = true
	}

	return func(item *itemv1.Item) bool {
		return statusMap[item.Status]
	}
}

func IDsFilter(ids []string) FilterFunc {
	if len(ids) == 0 {
		return nil
	}

	idMap := make(map[string]bool, len(ids))
	for _, id := range ids {
		idMap[id] = true
	}

	return func(item *itemv1.Item) bool {
		return idMap[item.Id]
	}
}

func ApplyItemFilter(protoFilter *itemv1.ItemFilter) FilterFunc {
	if protoFilter == nil {
		return func(*itemv1.Item) bool { return true }
	}

	builder := NewFilterBuilder()

	if protoFilter.Name != nil {
		builder.AddFilter(NameFilter(*protoFilter.Name))
	}

	if protoFilter.Description != nil {
		builder.AddFilter(DescriptionFilter(*protoFilter.Description))
	}

	builder.AddFilter(StatusFilter(protoFilter.Statuses))

	builder.AddFilter(IDsFilter(protoFilter.Ids))

	return builder.Build()
}

func ApplyItemFilters(protoFilters []*itemv1.ItemFilter) FilterFunc {
	if len(protoFilters) == 0 {
		return func(*itemv1.Item) bool { return true }
	}

	builder := NewFilterBuilder()

	for _, protoFilter := range protoFilters {
		builder.AddFilter(ApplyItemFilter(protoFilter))
	}

	return builder.Build()
}

func FilterItems(items []*itemv1.Item, filter FilterFunc) []*itemv1.Item {
	if filter == nil {
		return items
	}

	result := make([]*itemv1.Item, 0, len(items))
	for _, item := range items {
		if filter(item) {
			result = append(result, item)
		}
	}
	return result
}
