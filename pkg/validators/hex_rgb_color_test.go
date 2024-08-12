package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidHexRGBColor(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validHexRGBColor", ValidHexRGBColor)
	assert.Nil(t, err)

	color := "000000"
	err = validate.Var(color, "validHexRGBColor")
	assert.Nil(t, err)

	color = "000"
	err = validate.Var(color, "validHexRGBColor")
	assert.Nil(t, err)

	color = "e0e0e0"
	err = validate.Var(color, "validHexRGBColor")
	assert.Nil(t, err)

	color = "ffffff"
	err = validate.Var(color, "validHexRGBColor")
	assert.Nil(t, err)

	color = "FFFFFF"
	err = validate.Var(color, "validHexRGBColor")
	assert.Nil(t, err)
}

func TestInvalidHexRGBColor(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validHexRGBColor", ValidHexRGBColor)
	assert.Nil(t, err)

	color := "f"
	err = validate.Var(color, "validHexRGBColor")
	assert.NotNil(t, err)

	color = "fffffff"
	err = validate.Var(color, "validHexRGBColor")
	assert.NotNil(t, err)

	color = "gggggg"
	err = validate.Var(color, "validHexRGBColor")
	assert.NotNil(t, err)

	color = "#ffffff"
	err = validate.Var(color, "validHexRGBColor")
	assert.NotNil(t, err)
}
