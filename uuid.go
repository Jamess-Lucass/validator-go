package schema

import "github.com/google/uuid"

type UUIDSchema struct {
	Schema[uuid.UUID]
}

var _ ISchema = (*UUIDSchema)(nil)

func UUID() *UUIDSchema {
	return &UUIDSchema{}
}
