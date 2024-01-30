package schema

type ValidationResult struct {
	Errors []ValidationError
}

func (v *ValidationResult) IsValid() bool {
	return len(v.Errors) == 0
}

type ValidationError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

func (m *ValidationError) Error() string {
	return m.Message
}

type ISchema interface {
	Parse(value any) *ValidationResult
}

type Schema[T any] struct {
	validators []Validator[T]
}

type Validator[T any] struct {
	MessageFunc  func(T) string
	ValidateFunc func(T) bool
}
