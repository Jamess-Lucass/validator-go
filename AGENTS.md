# AGENTS.md

## Commands

```bash
go test -race ./...     # test (CI runs this across Go 1.18-1.26)
gofmt -s -w .           # format (CI checks this)
golangci-lint run ./... # lint (CI checks this)
```

## Architecture

- Root package `validator` (struct mode + engine) and subpackage `rule` (standalone rule values)
- Every rule attaches the same way, as a free function taking the validator and a field pointer (Go methods can't have type parameters, so there are no rule methods on `Validator` at all):
  - `validator.String(v, &s.Name).Min(2)`, `validator.Number(v, &s.Age).Gte(0)`, `validator.Slice(v, &s.Items)`, `validator.Map(v, &s.Labels)`, `validator.Bool`, `validator.Time`, `validator.UUID`, `validator.Complex`, `validator.Array`
  - `Number` covers all real numeric types; there are no Int/Float64 attach shorthands
- Struct flow: `v, err := validator.New(&s)`, add rules, then `v.Validate()` returns `*ValidationResult`
- Field names resolved via reflection from the pointer, preferring `json` struct tags
- Constraints use `~` type sets, so named types (`type ID string`, `type Age int`) work without casts
- Nil handling: nullable fields (`*string`) skip validation when nil unless `NotNil()` is chained. A nil pointer has no address, so its name can't be resolved by reflection; chain `WithName("...")` (on every rule type) to name it
- `New(ptr)` returns `(*Validator, error)`; the error is non-nil (and the validator nil) if ptr is not a non-nil pointer to a struct
- Standalone rule values live in `rule`: `rule.String()`, `rule.Int()`, `rule.Number[uint8]()`, `rule.Bool()`, `rule.Time()`, `rule.UUID()`, `rule.Complex128()`. They implement the `Rule` interface, coerce input (rejecting out-of-range values), and plug into `EachValue`, `EachKey`, and `Field`
- Unstructured data: `mv := validator.NewObject(m)`, then `mv.Field("key", rule.String().Email())`, then `mv.Validate()`. Nested objects use `validator.Object(func(omv *validator.ObjectValidator) { omv.Field(...) })`; `validator.MapOf`/`validator.SliceOf` apply one rule across all values/elements

## Testing

- Uses `testify/assert`
- Each type has a corresponding `_test.go` file
- CI tests against Go 1.18–1.26 with `-race`, plus golangci-lint and govulncheck jobs

## Style

- `golangci-lint` with `.golangci.yml` config (standard + gocritic, gosec, misspell, etc.)
- `gofmt -s` enforced
- Go 1.18 module; uses generics
