package util

func MustE[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func IgnoreE[T any](v T, err error) T {
	return v
}

func ZeroE[T any](v T, err error) T {
	if err != nil {
		var z T
		return z
	}
	return v
}
