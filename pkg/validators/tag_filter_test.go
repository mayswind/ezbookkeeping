package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestEmptyTagFilter(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validTagFilter", ValidTagFilter)
	assert.Nil(t, err)

	err = validate.Var("", "validTagFilter")
	assert.Nil(t, err)
}

func TestNoTag(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validTagFilter", ValidTagFilter)
	assert.Nil(t, err)

	err = validate.Var("none", "validTagFilter")
	assert.Nil(t, err)
}

func TestNoValidFilter(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validTagFilter", ValidTagFilter)
	assert.Nil(t, err)

	err = validate.Var(";", "validTagFilter")
	assert.NotNil(t, err)

	err = validate.Var(";;", "validTagFilter")
	assert.NotNil(t, err)
}

func TestValidOneFilterInTagFilters(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validTagFilter", ValidTagFilter)
	assert.Nil(t, err)

	err = validate.Var("0:1", "validTagFilter")
	assert.Nil(t, err)

	err = validate.Var("0:1,2,3", "validTagFilter")
	assert.Nil(t, err)

	err = validate.Var("1:1,2,3", "validTagFilter")
	assert.Nil(t, err)

	err = validate.Var("2:1,2,3", "validTagFilter")
	assert.Nil(t, err)

	err = validate.Var("3:1,2,3", "validTagFilter")
	assert.Nil(t, err)
}

func TestInvalidTagFilterType(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validTagFilter", ValidTagFilter)
	assert.Nil(t, err)

	err = validate.Var("a:1,2,3", "validTagFilter")
	assert.NotNil(t, err)

	err = validate.Var("-1:1,2,3", "validTagFilter")
	assert.NotNil(t, err)

	err = validate.Var("4:1,2,3", "validTagFilter")
	assert.NotNil(t, err)
}

func TestNoTagIdsInFilter(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validTagFilter", ValidTagFilter)
	assert.Nil(t, err)

	err = validate.Var("0", "validTagFilter")
	assert.NotNil(t, err)

	err = validate.Var("0:", "validTagFilter")
	assert.NotNil(t, err)
}

func TestInvalidTagIdsInFilter(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validTagFilter", ValidTagFilter)
	assert.Nil(t, err)

	err = validate.Var("0:abc", "validTagFilter")
	assert.NotNil(t, err)
}

func TestValidTwoFilterInTagFilters(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validTagFilter", ValidTagFilter)
	assert.Nil(t, err)

	err = validate.Var("0:1,2,3;2:4,5,6", "validTagFilter")
	assert.Nil(t, err)
}
