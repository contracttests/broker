package dsl

import "strings"

type PropertyPath string

func NewPropertyPath(propertyPath string) PropertyPath {
	return PropertyPath(propertyPath)
}

func (f *PropertyPath) Append(parts ...string) PropertyPath {
	separator := "."

	if string(*f) == "" {
		return PropertyPath(strings.Join(parts, separator))
	}

	return PropertyPath(strings.Join([]string{string(*f), strings.Join(parts, separator)}, separator))
}

func (f *PropertyPath) AppendArray() PropertyPath {
	return PropertyPath(f.String() + "[]")
}

func (f *PropertyPath) String() string {
	return string(*f)
}
