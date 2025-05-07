package util

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// IfNotNull accepts a default argument with a list of arguments and returns
// first non-nil argument. If not found then returns default.
func Coalesce[T any](args ...*T) *T {
	var r *T
	for _, a := range args {
		if a != nil {
			r = a
			break
		}
	}

	return r
}

// IfNull accepts a default argument with a list of arguments and returns
// first non-nil argument. If not found then returns default.
func IfNull[T any](d *T, args ...*T) *T {
	for _, a := range args {
		if a != nil {
			return a
		}
	}

	return d
}

// IfZero accepts a default argument with a list of arguments and returns the
// first non-zero argument. If not found then returns default.
func IfZero[T comparable](d T, args ...T) T {
	var z T
	for _, a := range args {
		if a != z {
			return a
		}
	}
	return d
}

// FromPtrSafe turns zero if no not-nil value found.
func FromPtrSafe[T any](d T, args ...*T) T {
	for _, a := range args {
		if a != nil {
			return *a
		}
	}

	var z T
	return z
}

// First accepts a list of arguments and returns the first non-zero argument.
// If not found then returns zero value of T.
func First[T comparable](args ...T) T {
	var z T
	for _, a := range args {
		if a != z {
			return a
		}
	}
	return z
}

// Map accepts a slice and a function and returns a new slice.
func Map[A, B any](elems []A, fn func(A) B) []B {
	result := make([]B, len(elems))
	for i, a := range elems {
		result[i] = fn(a)
	}
	return result
}

// QueryDecoder represents the query decoder type.
type QueryDecoder map[string]string

// Decode decodes the query string.
func (qd *QueryDecoder) Decode(val string) error {
	res := map[string]string{}
	v, err := url.ParseQuery(val)
	if err != nil {
		return err
	}

	for key, val := range v {
		if len(val) > 0 {
			res[key] = val[0]
		}
	}

	*qd = res
	return nil
}

// PrintJSON prints the json with indentation.
func PrintJSON(v any, p ...bool) {
	if len(p) > 0 && p[0] {
		printPrettyJSON(v)
		return
	}

	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println(v)
		return
	}
	fmt.Println(string(b))
}

func printPrettyJSON(v any) {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		fmt.Println(v)
		return
	}
	fmt.Println(string(b))
}

// StrBool converts string to boolean.
func StrBool(v string) bool {
	return strings.EqualFold(v, "true") || v == "1"
}
