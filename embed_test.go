package validator_test

import (
	"testing"
	"time"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

type embedBase struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type embedUser struct {
	embedBase
	Name string `json:"name"`
}

func TestEmbeddedStruct_FieldNames(t *testing.T) {
	user := embedUser{
		embedBase: embedBase{ID: "", CreatedAt: time.Time{}},
		Name:      "",
	}

	v, _ := validator.New(&user)
	v.String(&user.ID).NotEmpty()
	v.String(&user.Name).NotEmpty()
	v.Time(&user.CreatedAt).NotEmpty()
	result := v.Validate()

	assert.Len(t, result.Errors, 3)
	assert.Equal(t, "id", result.Errors[0].Field)
	assert.Equal(t, "name", result.Errors[1].Field)
	assert.Equal(t, "created_at", result.Errors[2].Field)
}

func TestEmbeddedStruct_WithNested(t *testing.T) {
	type innerAddr struct {
		City string `json:"city"`
	}
	type innerBase struct {
		ID string `json:"id"`
	}
	type innerUser struct {
		innerBase
		Name    string    `json:"name"`
		Address innerAddr `json:"address"`
	}

	user := innerUser{
		innerBase: innerBase{ID: ""},
		Name:      "",
		Address:   innerAddr{City: ""},
	}

	v, _ := validator.New(&user)
	v.String(&user.ID).NotEmpty()
	v.String(&user.Name).NotEmpty()
	v.String(&user.Address.City).NotEmpty()
	result := v.Validate()

	assert.Len(t, result.Errors, 3)
	assert.Equal(t, "id", result.Errors[0].Field)
	assert.Equal(t, "name", result.Errors[1].Field)
	assert.Equal(t, "address.city", result.Errors[2].Field)
}
