package main

import (
	"fmt"

	schema "github.com/Jamess-Lucass/validator-go"
)

type User struct {
	Firstname string
	Lastname  string
	Age       int
}

func main() {
	user := User{
		Firstname: "john",
		Lastname:  "doe",
		Age:       10,
	}

	s1 := schema.Object(map[string]schema.ISchema{
		"Firstname": schema.String().Min(4),
		"Lastname":  schema.String().Min(3),
		"Age":       schema.Int().LessThan(10),
	}).Parse(user)

	fmt.Printf("(1): is valid: %t\n", s1)
}
