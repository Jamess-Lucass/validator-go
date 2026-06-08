# AGENTS.md

## Commands

```bash
make test               # go test -v ./...
gofmt -s -w .           # format (CI checks this)
golangci-lint run ./... # lint (CI checks this)
```

## Architecture

- Single flat package `validator` (no cmd/, no internal/)
- FluentValidation-style builder API:
  - Struct: `v, err := validator.New(&s)` → add rules → `v.Validate()` returns `*ValidationResult`
  - Map: `v := validator.NewMap(m)` → `v.Field(key, rule)` → `v.Validate()` returns `*ValidationResult`
- Pointer-based field selection: `v.String(&s.Name).Min(2)` — field names resolved via reflection, preferring `json` struct tags
- Typed rules: String, Int, Float64, Bool, Time, UUID (`github.com/google/uuid`), Slice; numeric `Validate` (map mode) accepts any numeric kind
- Nil handling: nullable fields (`*string`) skip validation when nil unless `NotNil()` is chained. A nil pointer has no address, so its name can't be resolved by reflection — chain `WithName("...")` (on every rule type) to name it
- `New(ptr)` returns `(*Validator, error)`; the error is non-nil (and the validator nil) if ptr is not a non-nil pointer to a struct
- Slice validation uses a free function due to Go generics limitation: `validator.Slice(v, &s.Items).Each(...)`
- Map validation (unstructured data): `v := validator.NewMap(m)` → `v.Field("key", validator.String().Email())` → `v.Validate()`
- Nested maps use callback: `validator.Map(func(mv *validator.MapV) { mv.Field(...) })`
- `Rule` interface is implemented by all rule types for standalone use (maps, `EachValue`)
- `old-version/` contains the v1 code (zod-style schema API) for reference

## Testing

- Uses `testify/assert`
- Each type has a corresponding `_test.go` file
- CI tests against Go 1.18–1.23

## Style

- `golangci-lint` with `.golangci.yml` config (standard + gocritic, gosec, misspell, etc.)
- `gofmt -s` enforced
- Go 1.18 module; uses generics
