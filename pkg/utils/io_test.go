package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetImageContentType(t *testing.T) {
	fileName := "gif"
	expectedContentType := "image/gif"
	actualContentType := GetImageContentType(fileName)
	assert.Equal(t, expectedContentType, actualContentType)

	fileName = "bmp"
	expectedContentType = ""
	actualContentType = GetImageContentType(fileName)
	assert.Equal(t, expectedContentType, actualContentType)
}

func TestGetFileNameWithoutExtension(t *testing.T) {
	fileName := "name.ext"
	expectedName := "name"
	actualName := GetFileNameWithoutExtension(fileName)
	assert.Equal(t, expectedName, actualName)

	fileName = "C:\\name.ext"
	expectedName = "name"
	actualName = GetFileNameWithoutExtension(fileName)
	assert.Equal(t, expectedName, actualName)

	fileName = "/root/name.ext"
	expectedName = "name"
	actualName = GetFileNameWithoutExtension(fileName)
	assert.Equal(t, expectedName, actualName)

	fileName = "name"
	expectedName = "name"
	actualName = GetFileNameWithoutExtension(fileName)
	assert.Equal(t, expectedName, actualName)

	fileName = "C:\\name"
	expectedName = "name"
	actualName = GetFileNameWithoutExtension(fileName)
	assert.Equal(t, expectedName, actualName)

	fileName = "/root/name"
	expectedName = "name"
	actualName = GetFileNameWithoutExtension(fileName)
	assert.Equal(t, expectedName, actualName)

	fileName = "name."
	expectedName = "name"
	actualName = GetFileNameWithoutExtension(fileName)
	assert.Equal(t, expectedName, actualName)

	fileName = "C:\\name."
	expectedName = "name"
	actualName = GetFileNameWithoutExtension(fileName)
	assert.Equal(t, expectedName, actualName)

	fileName = "/root/name."
	expectedName = "name"
	actualName = GetFileNameWithoutExtension(fileName)
	assert.Equal(t, expectedName, actualName)

	fileName = ".ext"
	expectedName = ""
	actualName = GetFileNameWithoutExtension(fileName)
	assert.Equal(t, expectedName, actualName)

	fileName = "C:\\.ext"
	expectedName = ""
	actualName = GetFileNameWithoutExtension(fileName)
	assert.Equal(t, expectedName, actualName)

	fileName = "/root/.ext"
	expectedName = ""
	actualName = GetFileNameWithoutExtension(fileName)
	assert.Equal(t, expectedName, actualName)
}

func TestGetFileNameExtension(t *testing.T) {
	fileName := "name.ext"
	expectedExt := "ext"
	actualExt := GetFileNameExtension(fileName)
	assert.Equal(t, expectedExt, actualExt)

	fileName = "C:\\name.ext"
	expectedExt = "ext"
	actualExt = GetFileNameExtension(fileName)
	assert.Equal(t, expectedExt, actualExt)

	fileName = "/root/name.ext"
	expectedExt = "ext"
	actualExt = GetFileNameExtension(fileName)
	assert.Equal(t, expectedExt, actualExt)

	fileName = "name"
	expectedExt = ""
	actualExt = GetFileNameExtension(fileName)
	assert.Equal(t, expectedExt, actualExt)

	fileName = "C:\\name"
	expectedExt = ""
	actualExt = GetFileNameExtension(fileName)
	assert.Equal(t, expectedExt, actualExt)

	fileName = "/root/name"
	expectedExt = ""
	actualExt = GetFileNameExtension(fileName)
	assert.Equal(t, expectedExt, actualExt)

	fileName = "name."
	expectedExt = ""
	actualExt = GetFileNameExtension(fileName)
	assert.Equal(t, expectedExt, actualExt)

	fileName = "C:\\name."
	expectedExt = ""
	actualExt = GetFileNameExtension(fileName)
	assert.Equal(t, expectedExt, actualExt)

	fileName = "/root/name."
	expectedExt = ""
	actualExt = GetFileNameExtension(fileName)
	assert.Equal(t, expectedExt, actualExt)

	fileName = ".ext"
	expectedExt = "ext"
	actualExt = GetFileNameExtension(fileName)
	assert.Equal(t, expectedExt, actualExt)

	fileName = "C:\\.ext"
	expectedExt = "ext"
	actualExt = GetFileNameExtension(fileName)
	assert.Equal(t, expectedExt, actualExt)

	fileName = "/root/.ext"
	expectedExt = "ext"
	actualExt = GetFileNameExtension(fileName)
	assert.Equal(t, expectedExt, actualExt)
}
