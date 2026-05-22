package model

import (
	"sort"
	"strings"
)

func (r *CompatibilityReport) Canonical() []BreakingChange {
	out := append([]BreakingChange(nil), r.breaks...)

	sort.Slice(out, func(l, r int) bool {
		left, right := out[l], out[r]
		return left.CanonicalKey() < right.CanonicalKey()
	})

	return out
}

func (b BreakingChange) CanonicalKey() string {
	return strings.Join([]string{
		b.ContractInfo.Name,
		b.Resource.Endpoint,
		b.Resource.Method,
		b.Resource.StatusCode,
		b.Property,
		string(b.Reason),
	}, ";;")
}
