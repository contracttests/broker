package contract_differ

import "github.com/contracttesting/broker/internal/model"

type ChangeKind string

const (
	ChangeAdded    ChangeKind = "added"
	ChangeModified ChangeKind = "modified"
	ChangeRemoved  ChangeKind = "removed"
)

type ContractDiff struct {
	Resources map[string]ResourceChange
}

type ResourceChange struct {
	Kind       ChangeKind
	Resource   model.Resource
	Properties map[string]PropertyChange
}

type PropertyChange struct {
	Kind   ChangeKind
	Before model.Property
	After  model.Property
}

func Diff(prev *model.Contract, next *model.Contract) ContractDiff {
	var upcoming map[string]model.Resource

	if next != nil {
		upcoming = next.Resources
	}

	changes := map[string]ResourceChange{}

	for key, prevResource := range prev.Resources {
		nextResource, ok := upcoming[key]
		if !ok {
			changes[key] = removedResourceChange(prevResource)
			continue
		}
		if modified, ok := modifiedResourceChange(prevResource, nextResource); ok {
			changes[key] = modified
		}
	}

	for key, nextResource := range upcoming {
		if _, ok := prev.Resources[key]; !ok {
			changes[key] = addedResourceChange(nextResource)
		}
	}

	return ContractDiff{Resources: changes}
}

func addedResourceChange(r model.Resource) ResourceChange {
	changes := make(map[string]PropertyChange, len(r.Properties))

	for path, property := range r.Properties {
		changes[path] = PropertyChange{Kind: ChangeAdded, After: property}
	}

	return ResourceChange{Kind: ChangeAdded, Resource: r, Properties: changes}
}

func removedResourceChange(r model.Resource) ResourceChange {
	changes := make(map[string]PropertyChange, len(r.Properties))

	for path, property := range r.Properties {
		changes[path] = PropertyChange{Kind: ChangeRemoved, Before: property}
	}

	return ResourceChange{Kind: ChangeRemoved, Resource: r, Properties: changes}
}

func modifiedResourceChange(prev, next model.Resource) (ResourceChange, bool) {
	changes := map[string]PropertyChange{}

	for path, prevProperty := range prev.Properties {
		nextProperty, ok := next.Properties[path]
		if !ok {
			changes[path] = PropertyChange{Kind: ChangeRemoved, Before: prevProperty}
			continue
		}
		if !prevProperty.IsSame(&nextProperty) {
			changes[path] = PropertyChange{Kind: ChangeModified, Before: prevProperty, After: nextProperty}
		}
	}

	for path, nextProperty := range next.Properties {
		if _, ok := prev.Properties[path]; !ok {
			changes[path] = PropertyChange{Kind: ChangeAdded, After: nextProperty}
		}
	}

	if len(changes) == 0 {
		return ResourceChange{}, false
	}

	return ResourceChange{Kind: ChangeModified, Resource: next, Properties: changes}, true
}
