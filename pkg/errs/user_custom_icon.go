package errs

import "net/http"

// Error codes related to user custom icons
var (
	ErrUserCustomIconIdInvalid         = NewNormalError(NormalSubcategoryUserCustomIcon, 0, http.StatusBadRequest, "user custom icon id is invalid")
	ErrUserCustomIconNotFound          = NewNormalError(NormalSubcategoryUserCustomIcon, 1, http.StatusBadRequest, "user custom icon not found")
	ErrNoUserCustomIcon                = NewNormalError(NormalSubcategoryUserCustomIcon, 2, http.StatusBadRequest, "no user custom icon")
	ErrUserCustomIconIsEmpty           = NewNormalError(NormalSubcategoryUserCustomIcon, 3, http.StatusBadRequest, "user custom icon is empty")
	ErrUserCustomIconeNotExists        = NewNormalError(NormalSubcategoryUserCustomIcon, 4, http.StatusNotFound, "user custom icon not exists")
	ErrUserCustomIconExtensionInvalid  = NewNormalError(NormalSubcategoryUserCustomIcon, 5, http.StatusBadRequest, "user custom icon file extension invalid")
	ErrExceedMaxUserCustomIconFileSize = NewNormalError(NormalSubcategoryUserCustomIcon, 6, http.StatusBadRequest, "exceed the maximum size of user custom icon file")
	ErrUserCustomIconDimensionsInvalid = NewNormalError(NormalSubcategoryUserCustomIcon, 7, http.StatusBadRequest, "user custom icon dimensions must not exceed 256 pixels")
	ErrUserCustomIconInUse             = NewNormalError(NormalSubcategoryUserCustomIcon, 8, http.StatusBadRequest, "user custom icon is in use")
)
