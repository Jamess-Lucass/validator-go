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

	// 1
	s1 := schema.Object(map[string]schema.ISchema{
		"Firstname": schema.String().Min(5),
		"Lastname":  schema.Int(),
		"Age":       schema.Int().Lt(10),
	}).Parse(user)

	fmt.Printf("(1): is valid: %t\n", s1.IsValid())

	// 2
	s2 := schema.String().Refine(func(value string) bool {
		return value == "john"
	}).Parse("john")

	fmt.Printf("(2): is valid: %t\n", s2.IsValid())

	// 3
	s3 := schema.Object(map[string]schema.ISchema{
		"Firstname": schema.String(),
	}).Refine(func(value map[string]interface{}) bool {
		return value["Firstname"] == "john"
	}).Parse(user)

	fmt.Printf("(3): is valid: %t\n", s3.IsValid())
}
