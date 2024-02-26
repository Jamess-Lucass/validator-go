<p align="center">
  <h1 align="center">Validator-go</h1>
  <p align="center">
    <br/>
    <a href="https://github.com/colinhacks/zod">Zod</a> inspired schema validation in Go.
  </p>
</p>
<br/>
<p align="center">
<a href="https://github.com/Jamess-Lucass/validator-go/actions?query=branch%3Amain"><img src="https://github.com/Jamess-Lucass/validator-go/actions/workflows/test.yml/badge.svg?event=push&branch=main" alt="CI Test Status" /></a>
<a href="https://opensource.org/licenses/MIT" rel="nofollow"><img src="https://img.shields.io/github/license/Jamess-Lucass/validator-go" alt="License"></a>
</p>

## Table of contents

- [Table of contents](#table-of-contents)
- [Introduction](#introduction)
- [Installation](#installation)
- [Basic usage](#basic-usage)
- [Primitives](#primitives)
- [Literals](#literals)
- [Strings](#strings)
- [Ints](#ints)
- [Booleans](#booleans)
- [Objects](#objects)
- [Schema methods](#schema-methods)
  - [`.parse`](#parse)
  - [`.refine`](#refine)

## Introduction

validator-go is a simple and extensible validation library for Go. It's heavily inspired by [Zod](https://github.com/colinhacks/zod) and shares a lot of the same interfaces. This library provides a set of validation schemas for different data types, such as integers, strings, and booleans, and allows you to refine these schemas with custom validation rules.

## Installation

Ensure you have Go installed ([download](https://go.dev/dl/)). Version `1.21` or higher is required.

```bash
go get -u github.com/Jamess-Lucass/validator-go
```

## Basic usage

Creating a string schema

```go
import (
    schema "github.com/Jamess-Lucass/validator-go"
)

// Creating a schema for strings
mySchema := schema.String();

// Parsing
mySchema.Parse("john"); // => *schema.ValidationResult
mySchema.Parse("john").IsValid(); // => true
mySchema.Parse(12).Errors; // => []

mySchema.Parse(12).IsValid(); // => false
mySchema.Parse(12).Errors; // => [{ "path": "", "message": "Expected string, received int" }]
```

Creating an object schema

```go
import (
    schema "github.com/Jamess-Lucass/validator-go"
)

// Creating a schema for an object
mySchema := schema.Object(map[string]schema.ISchema{
    "Username": schema.String().Min(5),
})

// Parsing
type User struct {
	Username string
}

user1 := User{
    Username: "john_doe",
}

mySchema.Parse(user1).IsValid(); // => true
mySchema.Parse(user1).Errors; // => []

user2 := User{
    Username: "john",
}
mySchema.Parse(user2).IsValid(); // => false
mySchema.Parse(user2).Errors; // => [{ "path": "Firstname","message": "String must contain at least 5 character(s)" }]
```

## Primitives

```go
schema.String()
schema.Int()
schema.Bool()
```

## Literals

```go
john := schema.Literal("john");
four := schema.Literal(4);
trueSchema := schema.Literal(true);
```

## Strings

```go
schema.String().Max(2)
schema.String().Min(2)
schema.String().Length(2)
schema.String().Url()
schema.String().Includes(string)
schema.String().StartsWith(string)
schema.String().EndsWith(string)
```

## Ints

```go
schema.Int().Lt(2)
schema.Int().Lte(2)
schema.Int().Gt(2)
schema.Int().Gte(2)

schema.Int().Positive() // > 0
schema.Int().Nonnegative() // >= 0
schema.Int().Negative() // < 0
schema.Int().Nonpositive() // <= 0

schema.Int().MultipleOf(2)
```

## Booleans

```go
schema.Bool()
```

## Objects

```go
mySchema := schema.Object(map[string]schema.ISchema{
    "Username": schema.String().Min(5),
    "Firstname": schema.String().Min(2).Max(128),
    "Age": schema.Int().Gte(18),
    "IsVerified": schema.Bool(),
})
```

## Schema methods

All schemas contain certain methods.

### `.Parse`

Given any schema, you may call the `.Parse` method and pass through any data to check it's validity against the schema.

```go
mySchema := schema.String()

mySchema.Parse("john"); // => *schema.ValidationResult

mySchema.Parse("john").IsValid(); // => true
mySchema.Parse(2).Errors; // => []

mySchema.Parse(2).IsValid(); // => false
mySchema.Parse(2).Errors; // => [{ "path": "", "message": "Expected string, received int" }]
```

### `.Refine`

You may provide custom validation logic with the `.Refine` method. This method must return `true` or `false` to represent whether validation should be considered successful or unsuccessful.

```go
mySchema := schema.String().Refine(func(value string) bool {
    return value == "custom_value"
})
```

This is helpful if you need to perform some business-level validation. For example, checking a database for some value or making a HTTP request to assert something.

```go
verifiedUserSchema := schema.String().Refine(func(value string) bool {
    // Fetch user from data.
    // ensure `is_verified` field is true.
    return true
})
```

This may also be used in conjunction with `.Object`

```go
type User struct {
	Firstname string
	Lastname  string
	Age       int
}

user := User{
    Firstname: "john",
    Lastname:  "doe",
    Age:       10,
}

mySchema := schema.Object(map[string]schema.ISchema{
    "Firstname": schema.String().Refine(func(value string) bool {
        return value == "john" || value == "doe"
    }),
    "Lastname": schema.Int().Refine(func(value int) bool {
        return value == 10
    }),
    "Age": schema.Int().Lt(10),
}).Refine(func(value map[string]interface{}) bool {
    if value["Firstname"] == "doe" {
        if age, ok := value["Age"].(int); ok {
            return age < 5
        }
    }

    return true
})
```
