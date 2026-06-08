<p align="center">
  <h1 align="center">validator-go</h1>
  <p align="center">FluentValidation-style struct validation for Go.</p>
</p>
<br/>
<p align="center">
<a href="https://github.com/Jamess-Lucass/validator-go/actions?query=branch%3Amain"><img src="https://github.com/Jamess-Lucass/validator-go/actions/workflows/ci.yml/badge.svg?event=push&branch=main" alt="CI" /></a>
<a href="https://pkg.go.dev/github.com/Jamess-Lucass/validator-go"><img src="https://pkg.go.dev/badge/github.com/Jamess-Lucass/validator-go.svg" alt="Go Reference"></a>
<a href="https://opensource.org/licenses/Apache-2.0" rel="nofollow"><img src="https://img.shields.io/github/license/Jamess-Lucass/validator-go" alt="License"></a>
</p>

## Installation

```bash
go get github.com/Jamess-Lucass/validator-go
```

Requires Go 1.18+.

## Quick Start

Create a validator for a struct pointer, declare rules against its fields, then
call `Validate`:

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

user := User{Name: "J", Email: "bad", Age: -1}

v, err := validator.New(&user)
if err != nil {
    // &user must be a non-nil pointer to a struct
    return
}
v.String(&user.Name).NotEmpty().Min(2).Max(50)
v.String(&user.Email).NotEmpty().Email()
v.Int(&user.Age).Gte(0).Lte(150)

result := v.Validate()

result.IsValid() // false
result.Errors    // [{Field: "name", Message: "must be at least 2 characters"}, ...]
```

`New` returns an error if the argument is not a non-nil pointer to a struct.
Field names are resolved by reflection, preferring the `json` struct tag and
falling back to the Go field name. Rules read the field's value at `Validate`
time, so you can declare them up front and mutate the struct in between.

## String

```go
v.String(&s.Field).NotEmpty()
v.String(&s.Field).Min(n)
v.String(&s.Field).Max(n)
v.String(&s.Field).Length(min, max)
v.String(&s.Field).Email()
v.String(&s.Field).Url()
v.String(&s.Field).Includes("substr")
v.String(&s.Field).StartsWith("prefix")
v.String(&s.Field).EndsWith("suffix")
v.String(&s.Field).Must(func(val string) bool { return true })
v.String(&s.Field).WithMessage("custom error")
```

`Min`, `Max`, and `Length` count runes (Unicode code points), not bytes.

## Int

```go
v.Int(&s.Field).NotEmpty() // must not be zero
v.Int(&s.Field).Gt(n)
v.Int(&s.Field).Gte(n)
v.Int(&s.Field).Lt(n)
v.Int(&s.Field).Lte(n)
v.Int(&s.Field).Positive()
v.Int(&s.Field).Negative()
v.Int(&s.Field).Nonnegative()
v.Int(&s.Field).Nonpositive()
v.Int(&s.Field).MultipleOf(n)
v.Int(&s.Field).Must(func(val int) bool { return true })
```

## Float64

```go
v.Float64(&s.Field).NotEmpty() // must not be zero
v.Float64(&s.Field).Gt(n)
v.Float64(&s.Field).Gte(n)
v.Float64(&s.Field).Lt(n)
v.Float64(&s.Field).Lte(n)
v.Float64(&s.Field).Positive()
v.Float64(&s.Field).Negative()
v.Float64(&s.Field).Nonnegative()
v.Float64(&s.Field).Nonpositive()
v.Float64(&s.Field).MultipleOf(n)
v.Float64(&s.Field).Must(func(val float64) bool { return true })
```

## Bool

```go
v.Bool(&s.Field).NotEmpty() // must be true
v.Bool(&s.Field).Must(func(val bool) bool { return true })
```

## Time

```go
v.Time(&s.Field).NotEmpty() // must not be the zero time
v.Time(&s.Field).After(t)
v.Time(&s.Field).Before(t)
v.Time(&s.Field).Must(func(val time.Time) bool { return true })
```

## UUID

Validates `github.com/google/uuid` values.

```go
v.UUID(&s.ID).NotEmpty() // must not be the nil UUID (all zeros)
v.UUID(&s.ID).Must(func(val uuid.UUID) bool { return val.Version() == 4 })
v.UUID(&s.ID).WithMessage("custom error")
```

## Slice

`Slice` is a free function because Go methods cannot have type parameters.

```go
validator.Slice(v, &s.Tags).NotNil()  // fails on a nil slice (distinct from empty)
validator.Slice(v, &s.Tags).NotEmpty()
validator.Slice(v, &s.Tags).Min(n)
validator.Slice(v, &s.Tags).Max(n)
validator.Slice(v, &s.Tags).Length(min, max)
validator.Slice(v, &s.Tags).Must(func(val []string) bool { return true })

// Validate each primitive element
validator.Slice(v, &s.Tags).EachValue(validator.String().NotEmpty().Min(2))

// Validate each struct element
validator.Slice(v, &s.Items).Each(func(item *Item, sv *validator.Validator) {
    sv.String(&item.Name).NotEmpty()
    sv.Int(&item.Quantity).Gte(1)
})
```

## Nullable Fields

Pointer fields (`*string`, `*int`, `*time.Time`, ...) are skipped when nil. Pass
the pointer directly and use `NotNil()` to require a value.

```go
type Order struct {
    Notes       *string    `json:"notes"`
    ScheduledAt *time.Time `json:"scheduled_at"`
}

v, _ := validator.New(&order)
v.String(order.Notes).Min(1).Max(500)                       // skipped if nil
v.Time(order.ScheduledAt).NotNil().WithName("scheduled_at") // error if nil
```

A nil pointer has no address, so its field name cannot be resolved by
reflection — the error's `Field` falls back to `"unknown"`. Pair `NotNil()` with
`WithName("...")` (available on every rule type) to give such fields a stable
name in the error output.

## Nested & Embedded Structs

Nested struct fields resolve with dot-notation. Embedded struct fields resolve
without a prefix.

```go
type Address struct {
    City string `json:"city"`
}

type User struct {
    Base                       // embedded: fields resolve as "id", "created_at"
    Name    string  `json:"name"`
    Address Address `json:"address"` // nested: resolves as "address.city"
}

v, _ := validator.New(&user)
v.String(&user.Name).NotEmpty()
v.String(&user.Address.City).NotEmpty().Min(2) // error field: "address.city"
```

## Custom Validation

```go
v.String(&s.Code).Must(func(val string) bool {
    return strings.HasPrefix(val, "PRJ-")
}).WithMessage("must start with PRJ-")
```

`WithMessage` applies to the preceding rule.

## Cross-Field & Conditional Validation

Use plain Go control flow:

```go
v, _ := validator.New(&order)

if order.Express {
    v.Time(order.ScheduledAt).NotNil().WithName("scheduled_at").WithMessage("required for express orders")
}

if order.ScheduledAt != nil {
    v.Time(&order.ExpiresAt).Must(func(t time.Time) bool {
        return t.After(*order.ScheduledAt)
    }).WithMessage("must be after scheduled_at")
}

result := v.Validate()
```

## Map Validation

For unstructured data (`map[string]any`) such as webhooks or dynamic forms, use
`NewMap` and declare a rule per key. Numbers decoded from JSON arrive as
`float64`; the numeric rules accept any numeric kind.

```go
v := validator.NewMap(payload)
v.Field("email", validator.String().NotEmpty().Email())
v.Field("age", validator.Int().Gte(18))
v.Field("address", validator.Map(func(mv *validator.MapV) {
    mv.Field("city", validator.String().NotEmpty().Min(2))
    mv.Field("country", validator.String().NotEmpty())
}).NotNil())
v.Field("tags", validator.SliceOf(validator.String().NotEmpty()).Min(1).Max(10))

result := v.Validate()
```

The standalone rule constructors (`validator.String()`, `validator.Int()`, ...)
implement the `Rule` interface, so the same rules work in maps, in `EachValue`,
and in `SliceOf`.

## Error Format

```go
type ValidationError struct {
    Field   string `json:"field"`   // "name", "address.city", "items[0].price"
    Message string `json:"message"` // "must not be empty", "must be at least 2 characters"
}

type ValidationResult struct {
    Errors []ValidationError
}

result.IsValid() // true if no errors
```

## License

[Apache 2.0](LICENSE)
