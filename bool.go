package schema

type BoolSchema struct {
	Schema[bool]
}

var _ ISchema = (*BoolSchema)(nil)

func Bool() *BoolSchema {
	return &BoolSchema{}
}
