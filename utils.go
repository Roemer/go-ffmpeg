package goffmpeg

func Ternary[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func Ptr[T any](v T) *T {
	return &v
}
