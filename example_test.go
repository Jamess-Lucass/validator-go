package validator_test

import (
	"fmt"

	validator "github.com/Jamess-Lucass/validator-go"
)

func Example() {
	type User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	user := User{Name: "J", Email: "bad", Age: -1}

	v, err := validator.New(&user)
	if err != nil {
		panic(err)
	}
	v.String(&user.Name).NotEmpty().Min(2)
	v.String(&user.Email).Email()
	v.Int(&user.Age).Gte(0)

	result := v.Validate()

	fmt.Println(result.IsValid())
	for _, e := range result.Errors {
		fmt.Printf("%s: %s\n", e.Field, e.Message)
	}
	// Output:
	// false
	// name: must be at least 2 characters
	// email: must be a valid email address
	// age: must be greater than or equal to 0
}
