package schema

type ISchema interface {
	Parse(value any) bool
}

type Schema interface {
	ISchema
}
