package errs

import "net/http"

// Error codes related to model context protocol server
var (
	ErrMCPServerNotEnabled = NewNormalError(NormalSubcategoryModelContextProtocol, 0, http.StatusBadRequest, "mcp server is not enabled")
	ErrMCPToolNotEnabled   = NewNormalError(NormalSubcategoryModelContextProtocol, 1, http.StatusForbidden, "mcp tool is not enabled")
)
