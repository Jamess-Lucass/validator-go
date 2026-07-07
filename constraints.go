package validator

// Integer matches every built-in integer type and any named type whose
// underlying type is one of them (e.g. `type Age int`).
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Float matches the floating-point types and named types based on them.
type Float interface {
	~float32 | ~float64
}

// Real matches every ordered numeric type (integers and floats).
type Real interface {
	Integer | Float
}

// ComplexNumber matches the complex types and named types based on them.
// Complex values are not ordered, so only presence and custom rules apply.
type ComplexNumber interface {
	~complex64 | ~complex128
}
