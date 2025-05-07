package main

import (
	"strings"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func flattenMap(prefix string, data map[string]any, sep string, result map[string]any) {
	for k, v := range data {
		key := k
		if prefix != "" {
			key = prefix + sep + k
		}

		switch val := v.(type) {
		case map[string]any:
			// Recursively flatten nested maps
			flattenMap(key, val, sep, result)
		case nil:
			result[key] = nil
		default:
			// Store non-map values directly
			result[key] = val
		}
	}
}

// FlatMapConfig flattens the given map of config variables.
// ex: {"a": {"b": "c"}} => {"a_b": "c"}
func FlatMapConfig(vars map[string]any) pulumi.StringMap {
	flattened := make(map[string]any)
	flattenMap("", vars, "_", flattened)

	result := pulumi.StringMap{}
	for k, v := range flattened {
		val := pulumi.Sprintf("%v", v)
		k = strings.ToUpper(k)
		result[k] = val
	}
	return result
}

// MustSlice converts an interface to a slice of a specific type.
func MustSlice[T any](v any) []T {
	if v == nil {
		return nil
	}

	switch v := v.(type) {
	case []T:
		return v
	case []any:
		res := make([]T, len(v))
		for i, e := range v {
			if _, ok := e.(T); !ok {
				panic("unexpected type")
			}
			res[i] = e.(T)
		}
		return res
	default:
		panic("unexpected type")
	}
}

// MustSliceStr converts an interface to a slice of strings.
func MustSliceStr(v any) []string {
	return MustSlice[string](v)
}

// MustPSliceStr converts an interface to a pulumi.StringArray.
func MustPSliceStr(v any) pulumi.StringArray {
	return pulumi.ToStringArray(MustSliceStr(v))
}

// Select returns a if condA is true, otherwise b.
func Select[T any](condA bool, a, b T) T {
	if condA {
		return a
	}
	return b
}
