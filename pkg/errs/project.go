package errs

import "net/http"

// Project errors
var (
	ErrProjectIdInvalid            = NewNormalError(NormalSubcategoryProject, 0, http.StatusBadRequest, "project id is invalid")
	ErrProjectNotFound             = NewNormalError(NormalSubcategoryProject, 1, http.StatusBadRequest, "project not found")
	ErrProjectNameIsEmpty          = NewNormalError(NormalSubcategoryProject, 2, http.StatusBadRequest, "project name is empty")
	ErrProjectNameAlreadyExists    = NewNormalError(NormalSubcategoryProject, 3, http.StatusBadRequest, "project name already exists")
	ErrProjectInUseCannotBeDeleted = NewNormalError(NormalSubcategoryProject, 4, http.StatusBadRequest, "project is in use and cannot be deleted")
)
