package errs

import "net/http"

// Error codes related to transaction pictures
var (
	ErrTransactionPictureIdInvalid        = NewNormalError(NormalSubcategoryPicture, 0, http.StatusBadRequest, "transaction picture id is invalid")
	ErrTransactionPictureNotFound         = NewNormalError(NormalSubcategoryPicture, 1, http.StatusBadRequest, "transaction picture not found")
	ErrNoTransactionPicture               = NewNormalError(NormalSubcategoryPicture, 2, http.StatusBadRequest, "no transaction picture")
	ErrTransactionPictureIsEmpty          = NewNormalError(NormalSubcategoryPicture, 3, http.StatusBadRequest, "transaction picture is empty")
	ErrTransactionPictureNoExists         = NewNormalError(NormalSubcategoryPicture, 4, http.StatusNotFound, "transaction picture not exists")
	ErrTransactionPictureExtensionInvalid = NewNormalError(NormalSubcategoryPicture, 5, http.StatusNotFound, "transaction picture file extension invalid")
	ErrTransactionPictureTooMuch          = NewNormalError(NormalSubcategoryPicture, 6, http.StatusBadRequest, "transaction pictures too much")
)
