package rule_test

import (
	"testing"
	"time"

	"github.com/Jamess-Lucass/validator-go/rule"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestConstructors(t *testing.T) {
	assert.Empty(t, rule.String().NotEmpty().Validate("x"))
	assert.Empty(t, rule.Bool().NotEmpty().Validate(true))
	assert.Empty(t, rule.Int().Gte(18).Validate(21))
	assert.Empty(t, rule.Float64().Positive().Validate(1.5))
	assert.Empty(t, rule.Number[uint8]().Lte(200).Validate(100))
	assert.Empty(t, rule.Time().NotEmpty().Validate(time.Now()))
	assert.Empty(t, rule.UUID().NotEmpty().Validate(uuid.New()))
	assert.Empty(t, rule.Complex128().NotEmpty().Validate(1+2i))
	assert.Empty(t, rule.Complex64().NotEmpty().Validate(complex64(1)))

	errs := rule.String().Min(5).Validate("ab")
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be at least 5 characters", errs[0].Message)
}
