package main

import (
	"fmt"
	"strings"
	"time"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/Jamess-Lucass/validator-go/rule"
	"github.com/google/uuid"
)

type Name string

type StringExample struct {
	Name     Name   `json:"name"`
	Email    string `json:"email"`
	Website  string `json:"website"`
	Bio      string `json:"bio"`
	Filename string `json:"filename"`
	Code     string `json:"code"`
	Status   string `json:"status"`
	Greeting string `json:"greeting"`
}

type IntExample struct {
	Age         int `json:"age"`
	Score       int `json:"score"`
	Temperature int `json:"temperature"`
	BatchSize   int `json:"batch_size"`
	Balance     int `json:"balance"`
}

type Float64Example struct {
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount"`
	Rating      float64 `json:"rating"`
	Weight      float64 `json:"weight"`
	StepValue   float64 `json:"step_value"`
	Temperature float64 `json:"temperature"`
}

type BoolExample struct {
	Accepted *bool `json:"accepted"`
	Active   *bool `json:"active"`
	Verified bool  `json:"verified"`
}

type TimeExample struct {
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   time.Time  `json:"expires_at"`
	ScheduledAt *time.Time `json:"scheduled_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type SliceTag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SliceExample struct {
	Tags     []string   `json:"tags"`
	Scores   []int      `json:"scores"`
	Items    []SliceTag `json:"items"`
	Required []string   `json:"required"`
	Limited  []string   `json:"limited"`
}

type BaseEntity struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type NestedAddress struct {
	City   string `json:"city"`
	Street string `json:"street"`
}

type NestedProfile struct {
	Bio     string `json:"bio"`
	Website string `json:"website"`
}

type NestedExample struct {
	BaseEntity
	Name    string         `json:"name"`
	Address NestedAddress  `json:"address"`
	Profile *NestedProfile `json:"profile"`
}

type NullableExample struct {
	Name        string     `json:"name"`
	Nickname    *string    `json:"nickname"`
	Age         *int       `json:"age"`
	Score       *int       `json:"score"`
	ScheduledAt *time.Time `json:"scheduled_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type OrderItem struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

type OrderRequest struct {
	CustomerID  uuid.UUID   `json:"customer_id"`
	Items       []OrderItem `json:"items"`
	Notes       *string     `json:"notes"`
	Priority    int         `json:"priority"`
	Total       float64     `json:"total"`
	Express     bool        `json:"express"`
	ScheduledAt *time.Time  `json:"scheduled_at"`
	ExpiresAt   time.Time   `json:"expires_at"`
}

type Container struct {
	Image    string `json:"image"`
	Replicas int    `json:"replicas"`
}

type DeploymentExample struct {
	Name       string               `json:"name"`
	Labels     map[string]string    `json:"labels"`
	Containers map[string]Container `json:"containers"`
}

// Example 1: String
func example1StringValidation() *validator.ValidationResult {
	s := StringExample{
		Name:     "J",
		Email:    "not-an-email",
		Website:  "not-a-url",
		Bio:      "Hello world",
		Filename: "readme.txt",
		Code:     "ABC-123",
		Status:   "archived",
		Greeting: "",
	}

	v := must(validator.New(&s))
	// Name is a named string type; no cast needed.
	validator.String(v, &s.Name).NotEmpty().Min(2).Max(50)
	validator.String(v, &s.Name).Length(2, 50).WithMessage("name must be 2-50 characters")
	validator.String(v, &s.Email).NotEmpty().Email()
	validator.String(v, &s.Website).NotEmpty().URL()
	validator.String(v, &s.Bio).Includes("golang").WithMessage("bio must mention golang")
	validator.String(v, &s.Filename).StartsWith("doc_").WithMessage("filename must start with doc_")
	validator.String(v, &s.Code).Matches(`^[a-z]+-\d+$`).WithMessage("code must look like abc-123")
	validator.String(v, &s.Status).OneOf("active", "disabled")
	validator.String(v, &s.Greeting).NotEmpty().Must(func(val string) bool {
		return strings.HasPrefix(val, "Hello")
	}).WithMessage("greeting must start with Hello")
	return v.Validate()
}

// Example 2: Int
func example2IntValidation() *validator.ValidationResult {
	s := IntExample{
		Age:         -5,
		Score:       150,
		Temperature: 10,
		BatchSize:   7,
		Balance:     0,
	}

	v := must(validator.New(&s))
	validator.Number(v, &s.Balance).NotEmpty().WithMessage("balance must not be zero")
	validator.Number(v, &s.Age).Gte(0).WithMessage("age must be non-negative")
	validator.Number(v, &s.Balance).Gt(0).WithMessage("balance must be greater than zero")
	validator.Number(v, &s.Score).Lt(100).WithMessage("score must be less than 100")
	validator.Number(v, &s.Score).Lte(100).WithMessage("score must be at most 100")
	validator.Number(v, &s.Age).Positive().WithMessage("age must be positive")
	validator.Number(v, &s.Temperature).Negative().WithMessage("temperature must be negative")
	validator.Number(v, &s.Age).Nonnegative()
	validator.Number(v, &s.Temperature).Nonpositive().WithMessage("temperature must be non-positive")
	validator.Number(v, &s.BatchSize).MultipleOf(5).WithMessage("batch size must be a multiple of 5")
	validator.Number(v, &s.BatchSize).OneOf(10, 25, 50)
	validator.Number(v, &s.Score).Must(func(val int) bool {
		return val%2 == 0
	}).WithMessage("score must be even")
	return v.Validate()
}

// Example 3: Float64
func example3Float64Validation() *validator.ValidationResult {
	s := Float64Example{
		Price:       -9.99,
		Discount:    1.5,
		Rating:      6.0,
		Weight:      0,
		StepValue:   0.7,
		Temperature: 5.0,
	}

	v := must(validator.New(&s))
	validator.Number(v, &s.Weight).NotEmpty().WithMessage("weight must not be zero")
	validator.Number(v, &s.Price).Positive().WithMessage("price must be positive")
	validator.Number(v, &s.Temperature).Negative().WithMessage("temperature must be negative")
	validator.Number(v, &s.Weight).Gt(0).WithMessage("weight must be greater than zero")
	validator.Number(v, &s.Price).Gte(0).WithMessage("price must be non-negative")
	validator.Number(v, &s.Rating).Lt(5.0).WithMessage("rating must be less than 5.0")
	validator.Number(v, &s.Discount).Lte(1.0).WithMessage("discount must be at most 1.0 (100%)")
	validator.Number(v, &s.Price).Nonnegative()
	validator.Number(v, &s.Discount).Nonpositive().WithMessage("discount must be non-positive")
	validator.Number(v, &s.StepValue).MultipleOf(0.25).WithMessage("step must be a multiple of 0.25")
	validator.Number(v, &s.Rating).Must(func(val float64) bool {
		return val >= 1.0 && val <= 5.0
	}).WithMessage("rating must be between 1.0 and 5.0")
	return v.Validate()
}

// Example 4: Bool
func example4BoolValidation() *validator.ValidationResult {
	falseVal := false
	s := BoolExample{
		Accepted: nil,
		Active:   &falseVal,
		Verified: false,
	}

	v := must(validator.New(&s))
	validator.Bool(v, s.Accepted).NotNil().WithName("accepted").WithMessage("terms must be accepted")
	validator.Bool(v, s.Active).NotEmpty().WithMessage("account must be active")
	validator.Bool(v, &s.Verified).Must(func(val bool) bool {
		return val
	}).WithMessage("user must be verified")
	return v.Validate()
}

// Example 5: Time
func example5TimeValidation() *validator.ValidationResult {
	pastTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	s := TimeExample{
		CreatedAt:   time.Time{},
		ExpiresAt:   pastTime,
		ScheduledAt: nil,
		DeletedAt:   &pastTime,
	}

	v := must(validator.New(&s))
	now := time.Now()
	cutoff := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	validator.Time(v, &s.CreatedAt).NotEmpty().WithMessage("created_at is required")
	validator.Time(v, &s.ExpiresAt).After(now).WithMessage("expires_at must be in the future")
	validator.Time(v, &s.ExpiresAt).Before(cutoff).WithMessage("expires_at must be before 2021")
	validator.Time(v, s.ScheduledAt).NotNil().WithName("scheduled_at").WithMessage("scheduled_at is required")
	validator.Time(v, s.DeletedAt).Must(func(t time.Time) bool {
		return t.After(cutoff)
	}).WithMessage("deleted_at must be after 2021")
	return v.Validate()
}

// Example 6: Slice
func example6SliceValidation() *validator.ValidationResult {
	s := SliceExample{
		Tags:     []string{"go", "", "rust", ""},
		Scores:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		Items:    []SliceTag{{Name: "", Value: "x"}},
		Required: nil,
		Limited:  []string{"a", "b", "c", "d"},
	}

	v := must(validator.New(&s))
	validator.Slice(v, &s.Required).NotNil().WithMessage("required must not be nil")
	validator.Slice(v, &s.Required).NotEmpty().WithMessage("required must not be empty")
	validator.Slice(v, &s.Tags).Min(1).Max(3).WithMessage("tags must have at most 3 items")
	validator.Slice(v, &s.Limited).Length(1, 3).WithMessage("limited must have 1-3 items")
	validator.Slice(v, &s.Scores).Max(5).WithMessage("scores must have at most 5 items")

	validator.Slice(v, &s.Tags).EachValue(
		rule.String().NotEmpty().Min(2),
	)

	validator.Slice(v, &s.Items).NotEmpty().Each(func(item *SliceTag, sv *validator.Validator) {
		validator.String(sv, &item.Name).NotEmpty().Min(2).WithMessage("tag name must be at least 2 characters")
		validator.String(sv, &item.Value).NotEmpty()
	})

	validator.Slice(v, &s.Tags).Must(func(tags []string) bool {
		for _, t := range tags {
			if t == "go" {
				return true
			}
		}
		return false
	}).WithMessage("tags must contain 'go'")
	return v.Validate()
}

// Example 7: Nested Structs & Embedded Fields
func example7NestedStructs() *validator.ValidationResult {
	s := NestedExample{
		BaseEntity: BaseEntity{
			ID:        uuid.Nil,
			CreatedAt: time.Time{},
		},
		Name:    "",
		Address: NestedAddress{City: "A", Street: ""},
		Profile: &NestedProfile{Bio: "x", Website: "not-a-url"},
	}

	v := must(validator.New(&s))
	validator.Time(v, &s.CreatedAt).NotEmpty().WithMessage("created_at is required")
	validator.String(v, &s.Name).NotEmpty().Min(2).WithMessage("name must be at least 2 characters")
	validator.String(v, &s.Address.City).NotEmpty().Min(2).WithMessage("city must be at least 2 characters")
	validator.String(v, &s.Address.Street).NotEmpty().Min(5).WithMessage("street must be at least 5 characters")

	if s.Profile != nil {
		validator.String(v, &s.Profile.Bio).NotEmpty().Min(10).WithMessage("bio must be at least 10 characters")
		validator.String(v, &s.Profile.Website).NotEmpty().URL().WithMessage("website must be a valid URL")
	}

	validator.UUID(v, &s.ID).NotEmpty().WithMessage("id must not be a nil UUID")
	return v.Validate()
}

// Example 8: Nullable Fields & NotNil
func example8NullableFields() *validator.ValidationResult {
	zero := 0
	emptyStr := ""
	pastTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	s := NullableExample{
		Name:        "Alice",
		Nickname:    &emptyStr,
		Age:         nil,
		Score:       &zero,
		ScheduledAt: nil,
		DeletedAt:   &pastTime,
	}

	v := must(validator.New(&s))
	validator.String(v, &s.Name).NotEmpty().Min(2)
	validator.String(v, s.Nickname).NotEmpty().Min(2).WithMessage("nickname must be at least 2 characters")
	validator.Number(v, s.Age).Gte(0).Lte(150)
	validator.Number(v, s.Score).Positive().WithMessage("score must be positive")
	validator.Time(v, s.ScheduledAt).NotNil().WithName("scheduled_at").WithMessage("scheduled_at is required")
	validator.Time(v, s.DeletedAt).After(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)).
		WithMessage("deleted_at must be after 2022")
	return v.Validate()
}

// Example 9: Cross-Field & Conditional Validation
func example9CrossFieldConditional() *validator.ValidationResult {
	notes := ""
	order := OrderRequest{
		CustomerID:  uuid.New(),
		Priority:    6,
		Total:       -5.00,
		Notes:       &notes,
		Express:     true,
		ScheduledAt: nil,
		ExpiresAt:   time.Time{},
		Items: []OrderItem{
			{ProductName: "A", Quantity: 0, Price: -1.50},
			{ProductName: "", Quantity: 1001, Price: 0},
		},
	}

	v := must(validator.New(&order))
	// WithMessage binds to the rule right before it, so a range check names
	// each bound separately.
	validator.Number(v, &order.Priority).
		Gte(1).WithMessage("priority must be between 1 and 5").
		Lte(5).WithMessage("priority must be between 1 and 5")
	validator.Number(v, &order.Total).Positive().WithMessage("order total must be greater than zero")
	validator.String(v, order.Notes).NotEmpty().WithMessage("notes must not be empty if provided").Max(500)
	validator.Time(v, &order.ExpiresAt).NotEmpty().WithMessage("expires_at is required")

	if order.Express {
		validator.Time(v, order.ScheduledAt).NotNil().WithName("scheduled_at").WithMessage("scheduled_at is required for express orders")
	}

	if order.ScheduledAt != nil {
		validator.Time(v, &order.ExpiresAt).Must(func(t time.Time) bool {
			return t.After(*order.ScheduledAt)
		}).WithMessage("expires_at must be after scheduled_at")
	}

	validator.Slice(v, &order.Items).Min(1).Max(50).Each(func(item *OrderItem, sv *validator.Validator) {
		validator.String(sv, &item.ProductName).NotEmpty().Min(2).WithMessage("product name must be at least 2 characters")
		validator.Number(sv, &item.Quantity).
			Gte(1).WithMessage("quantity must be between 1 and 1000").
			Lte(1000).WithMessage("quantity must be between 1 and 1000")
		validator.Number(sv, &item.Price).Positive().WithMessage("price must be positive")
	})
	return v.Validate()
}

// Example 10: Map Validation
func example10MapValidation() *validator.ValidationResult {
	payload := map[string]any{
		"user_id":  "not-a-uuid",
		"email":    "not-an-email",
		"age":      float64(10),
		"score":    float64(3.7),
		"price":    float64(-5.5),
		"active":   "yes",
		"verified": false,
		"created":  "2020-01-01T00:00:00Z",
		"address": map[string]any{
			"city":    "N",
			"country": "",
		},
		"tags":   []any{"go", "", "rust"},
		"scores": []any{float64(1), float64(2)},
	}

	v := validator.NewObject(payload)
	v.Field("user_id", rule.UUID().NotNil())
	v.Field("email", rule.String().NotEmpty().Email())
	v.Field("age", rule.Int().Gte(18).Lte(120))
	v.Field("score", rule.Int().Gte(0))
	v.Field("price", rule.Float64().Positive().Gte(0))
	v.Field("active", rule.Bool().NotEmpty())
	v.Field("verified", rule.Bool().NotEmpty())
	v.Field("created", rule.Time().NotEmpty().After(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))
	v.Field("address", validator.Object(func(mv *validator.ObjectValidator) {
		mv.Field("city", rule.String().NotEmpty().Min(2))
		mv.Field("country", rule.String().NotEmpty())
	}))
	v.Field("profile", validator.Object(func(mv *validator.ObjectValidator) {
		mv.Field("bio", rule.String().NotEmpty())
	}).NotNil())
	v.Field("tags", validator.SliceOf(rule.String().NotEmpty()).NotEmpty().Min(1).Max(5))
	v.Field("scores", validator.SliceOf(rule.Int().Positive()).Min(3))
	v.Field("missing", rule.String().NotNil())
	return v.Validate()
}

// Example 11: Typed Map Fields
func example11TypedMapFields() *validator.ValidationResult {
	d := DeploymentExample{
		Name:   "api",
		Labels: map[string]string{"env": "", "team": "core", "": "x"},
		Containers: map[string]Container{
			"web": {Image: "", Replicas: 0},
			"db":  {Image: "postgres:16", Replicas: 1},
		},
	}

	v := must(validator.New(&d))
	validator.String(v, &d.Name).NotEmpty()

	validator.Map(v, &d.Labels).NotEmpty().HasKey("region").
		EachKey(rule.String().NotEmpty()).
		EachValue(rule.String().NotEmpty())

	validator.Map(v, &d.Containers).Min(1).Each(func(name string, c *Container, sv *validator.Validator) {
		validator.String(sv, &c.Image).NotEmpty()
		validator.Number(sv, &c.Replicas).Gte(1)
	})

	return v.Validate()
}

func main() {
	examples := []struct {
		name string
		fn   func() *validator.ValidationResult
	}{
		{"String Validation", example1StringValidation},
		{"Int Validation", example2IntValidation},
		{"Float64 Validation", example3Float64Validation},
		{"Bool Validation", example4BoolValidation},
		{"Time Validation", example5TimeValidation},
		{"Slice Validation", example6SliceValidation},
		{"Nested Structs & Embedded Fields", example7NestedStructs},
		{"Nullable Fields & NotNil", example8NullableFields},
		{"Cross-Field & Conditional Validation", example9CrossFieldConditional},
		{"Map Validation", example10MapValidation},
		{"Typed Map Fields", example11TypedMapFields},
	}

	for i, ex := range examples {
		fmt.Println(strings.Repeat("=", 60))
		fmt.Printf("EXAMPLE %d: %s\n", i+1, ex.name)
		fmt.Println(strings.Repeat("=", 60))

		printResult(ex.fn())

		fmt.Println()
	}
}

func printResult(result *validator.ValidationResult) {
	if result.IsValid() {
		fmt.Println("VALID")
		return
	}

	fmt.Printf("%d error(s):\n", len(result.Errors))

	for _, err := range result.Errors {
		if err.Field == "" {
			fmt.Printf("  - %-30s %s\n", "(root)", err.Message)
		} else {
			fmt.Printf("  - %-30s %s\n", err.Field, err.Message)
		}
	}
}

// must panics if New fails; real callers should handle the error instead.
func must(v *validator.Validator, err error) *validator.Validator {
	if err != nil {
		panic(err)
	}
	return v
}
