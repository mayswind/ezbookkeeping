package errs

import "net/http"

// Error codes related to import presets
var (
	ErrImportPresetIdInvalid    = NewNormalError(NormalSubcategoryImportPreset, 0, http.StatusBadRequest, "import preset id is invalid")
	ErrImportPresetNotFound     = NewNormalError(NormalSubcategoryImportPreset, 1, http.StatusBadRequest, "import preset not found")
	ErrImportPresetDataInvalid  = NewNormalError(NormalSubcategoryImportPreset, 2, http.StatusBadRequest, "import preset data is invalid")
	ErrImportPresetDataTooLarge = NewNormalError(NormalSubcategoryImportPreset, 3, http.StatusBadRequest, "import preset data is too large")
)
