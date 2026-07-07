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
call `Validate`. Every rule attaches the same way: a free function taking the
validator and a pointer to the field.

```go
import (
    validator "github.com/Jamess-Lucass/validator-go"
)

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
validator.String(v, &user.Name).NotEmpty().Min(2).Max(50)
validator.String(v, &user.Email).NotEmpty().Email()
validator.Number(v, &user.Age).Gte(0).Lte(150)

result := v.Validate()

result.IsValid() // false
result.Errors    // [{Field: "name", Message: "must be at least 2 characters"}, ...]
```

`New` returns an error if the argument is not a non-nil pointer to a struct.
Field names are resolved by reflection, preferring the `json` struct tag and
falling back to the Go field name. Rules read the field's value at `Validate`
time, so you can declare them up front and mutate the struct in between.

The `rule` subpackage provides the same rules as standalone values for
unstructured data and per-element validation; see [Unstructured
Objects](#unstructured-objects) below.

## String

Named string types (`type ID string`) work too.

```go
validator.String(v, &s.Field).NotEmpty()
validator.String(v, &s.Field).Min(n)
validator.String(v, &s.Field).Max(n)
validator.String(v, &s.Field).Length(min, max)
validator.String(v, &s.Field).Email()
validator.String(v, &s.Field).URL()
validator.String(v, &s.Field).Matches(`^[a-z]+-\d+$`)
validator.String(v, &s.Field).Includes("substr")
validator.String(v, &s.Field).StartsWith("prefix")
validator.String(v, &s.Field).EndsWith("suffix")
validator.String(v, &s.Field).OneOf("active", "disabled")
validator.String(v, &s.Field).Must(func(val string) bool { return true })
validator.String(v, &s.Field).WithMessage("custom error")
```

`Min`, `Max`, and `Length` count runes (Unicode code points), not bytes.
`Matches` compiles the pattern once when the rule is declared; like
`regexp.MustCompile`, an invalid pattern panics.

## Numbers

`Number` works on any real numeric type, from `int` and `float64` to named
types like `type Age int`. The type is inferred from the field pointer.

```go
validator.Number(v, &s.Field).NotEmpty() // must not be zero
validator.Number(v, &s.Field).Gt(n)
validator.Number(v, &s.Field).Gte(n)
validator.Number(v, &s.Field).Lt(n)
validator.Number(v, &s.Field).Lte(n)
validator.Number(v, &s.Field).Positive()
validator.Number(v, &s.Field).Negative()
validator.Number(v, &s.Field).Nonnegative()
validator.Number(v, &s.Field).Nonpositive()
validator.Number(v, &s.Field).OneOf(1, 2, 3)
validator.Number(v, &s.Field).MultipleOf(n)
validator.Number(v, &s.Field).Finite() // rejects NaN and ±Inf
validator.Number(v, &s.Field).Must(func(val int) bool { return true })
```

Integer `MultipleOf` is exact; float `MultipleOf` uses a magnitude-relative
tolerance. Complex numbers aren't ordered, so they get only presence and
custom rules:

```go
validator.Complex(v, &s.Amplitude).NotEmpty() // must be != 0
```

## Bool

Named bool types work too.

```go
validator.Bool(v, &s.Field).NotEmpty() // must be true
validator.Bool(v, &s.Field).Must(func(val bool) bool { return true })
```

## Time

```go
validator.Time(v, &s.Field).NotEmpty() // must not be the zero time
validator.Time(v, &s.Field).After(t)
validator.Time(v, &s.Field).Before(t)
validator.Time(v, &s.Field).Must(func(val time.Time) bool { return true })
```

## UUID

Validates `github.com/google/uuid` values.

```go
validator.UUID(v, &s.ID).NotEmpty() // must not be the nil UUID (all zeros)
validator.UUID(v, &s.ID).Must(func(val uuid.UUID) bool { return val.Version() == 4 })
validator.UUID(v, &s.ID).WithMessage("custom error")
```

## Slice

Named slice types (`type Tags []string`) work too.

```go
validator.Slice(v, &s.Tags).NotNil()  // fails on a nil slice (distinct from empty)
validator.Slice(v, &s.Tags).NotEmpty()
validator.Slice(v, &s.Tags).Min(n)
validator.Slice(v, &s.Tags).Max(n)
validator.Slice(v, &s.Tags).Length(min, max)
validator.Slice(v, &s.Tags).Must(func(val []string) bool { return true })

// Validate each primitive element
validator.Slice(v, &s.Tags).EachValue(rule.String().NotEmpty().Min(2))

// Validate each struct element
validator.Slice(v, &s.Items).Each(func(item *Item, sv *validator.Validator) {
    validator.String(sv, &item.Name).NotEmpty()
    validator.Number(sv, &item.Quantity).Gte(1)
})
```

## Arrays

Slices use the typed `Slice` above. Fixed-size arrays can't be expressed with Go
generics (the length is part of the type), so `Array` validates them by reflection:

```go
type Matrix struct {
    Row [3]int `json:"row"`
}

validator.Array(v, &m.Row).Min(1).EachValue(rule.Int().Positive())
```

Array values are also accepted anywhere a slice is (maps, `SliceOf`, `EachValue`).

## Maps

`Map` validates a typed map field, such as a set of labels or HTTP headers.
Named map types (`type Labels map[string]string`) work too. Entries are
checked in sorted key order so errors always come out in the same order, even
though Go randomizes map iteration. Pointer keys are shown by the value they
point to, since an address would change from run to run.

```go
type Deployment struct {
    Name       string               `json:"name"`
    Labels     map[string]string    `json:"labels"`
    Containers map[string]Container `json:"containers"`
}

validator.Map(v, &d.Labels).NotNil()   // fails on a nil map (distinct from empty)
validator.Map(v, &d.Labels).NotEmpty()
validator.Map(v, &d.Labels).Min(n)
validator.Map(v, &d.Labels).Max(n)
validator.Map(v, &d.Labels).Length(min, max)
validator.Map(v, &d.Labels).HasKey("env")
validator.Map(v, &d.Labels).Must(func(m map[string]string) bool { return true })

// Validate every value: "labels[env]: must not be empty"
validator.Map(v, &d.Labels).EachValue(rule.String().NotEmpty())

// Validate every key: "labels[e]: key must be at least 2 characters"
validator.Map(v, &d.Labels).EachKey(rule.String().Min(2))

// Validate struct values, with errors like "containers[web].image: must not be empty"
validator.Map(v, &d.Containers).Each(func(name string, c *Container, sv *validator.Validator) {
    validator.String(sv, &c.Image).NotEmpty()
    validator.Number(sv, &c.Replicas).Gte(1)
})
```

Go map values aren't addressable, so `Each` passes a pointer to a copy of the
value.

## Nullable Fields

Pointer fields (`*string`, `*int`, `*time.Time`, ...) are skipped when nil. Pass
the pointer directly and use `NotNil()` to require a value.

```go
type Order struct {
    Notes       *string    `json:"notes"`
    ScheduledAt *time.Time `json:"scheduled_at"`
}

v, _ := validator.New(&order)
validator.String(v, order.Notes).Min(1).Max(500)                       // skipped if nil
validator.Time(v, order.ScheduledAt).NotNil().WithName("scheduled_at") // error if nil
```

A nil pointer has no address, so its field name cannot be resolved by
reflection and the error's `Field` falls back to `"unknown"`. Chain
`WithName("...")` (available on every rule type) to give such fields a stable
name in the error output. The same applies when two pointer fields share one
target: the address alone cannot tell them apart, so name the rules
explicitly.

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
validator.String(v, &user.Name).NotEmpty()
validator.String(v, &user.Address.City).NotEmpty().Min(2) // error field: "address.city"
```

## Custom Validation

```go
validator.String(v, &s.Code).Must(func(val string) bool {
    return strings.HasPrefix(val, "PRJ-")
}).WithMessage("must start with PRJ-")
```

`WithMessage` applies to the preceding rule. Directly after `Each`,
`EachValue`, or `EachKey` it replaces the message of every error they
produce.

## Cross-Field & Conditional Validation

Use plain Go control flow:

```go
v, _ := validator.New(&order)

if order.Express {
    validator.Time(v, order.ScheduledAt).NotNil().WithName("scheduled_at").WithMessage("required for express orders")
}

if order.ScheduledAt != nil {
    validator.Time(v, &order.ExpiresAt).Must(func(t time.Time) bool {
        return t.After(*order.ScheduledAt)
    }).WithMessage("must be after scheduled_at")
}

result := v.Validate()
```

## Unstructured Objects

For unstructured data (`map[string]any`) such as webhooks or dynamic forms, use
`NewObject` and declare one rule per key. Unlike `New` it does not return an
error: any map is usable, including a nil one. A missing key validates as nil,
so it is skipped unless the rule chains `NotNil()`.

The standalone rule values come from the `rule` subpackage. Struct fields are
typed at compile time and validated strictly, but map values arrive however
the decoder produced them, so these rules coerce: numbers decoded from JSON as
`float64` and numeric strings both work. Error names come from the map key or
element index; `WithName` has no effect on standalone rules.

```go
import "github.com/Jamess-Lucass/validator-go/rule"

mv := validator.NewObject(payload)
mv.Field("email", rule.String().NotEmpty().Email())
mv.Field("age", rule.Int().Gte(18))
mv.Field("port", rule.Number[uint16]().Gt(0)) // coercion target of any width
mv.Field("address", validator.Object(func(omv *validator.ObjectValidator) {
    omv.Field("city", rule.String().NotEmpty().Min(2))
    omv.Field("country", rule.String().NotEmpty())
}).NotNil())
mv.Field("tags", validator.SliceOf(rule.String().NotEmpty()).Min(1).Max(10))
mv.Field("labels", validator.MapOf(rule.String().NotEmpty()).Max(20))

result := mv.Validate()
```

`Object` declares rules for known keys of a nested object, `MapOf` runs one
rule over every value, and `SliceOf` over every element. The same rule values
plug into `EachValue` and `EachKey` in struct mode. Values that don't fit the
target type, like a `uint64` above `MaxInt64` or a fractional float for an
integer rule, are rejected instead of wrapped.

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
