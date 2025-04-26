package duplicatechecker

// DuplicateCheckerType represents duplicate checker type
type DuplicateCheckerType uint8

// Types of uuid
const (
	DUPLICATE_CHECKER_TYPE_BACKGROUND_CRON_JOB DuplicateCheckerType = 0
	DUPLICATE_CHECKER_TYPE_NEW_ACCOUNT         DuplicateCheckerType = 1
	DUPLICATE_CHECKER_TYPE_NEW_SUBACCOUNT      DuplicateCheckerType = 2
	DUPLICATE_CHECKER_TYPE_NEW_CATEGORY        DuplicateCheckerType = 3
	DUPLICATE_CHECKER_TYPE_NEW_TRANSACTION     DuplicateCheckerType = 4
	DUPLICATE_CHECKER_TYPE_NEW_TEMPLATE        DuplicateCheckerType = 5
	DUPLICATE_CHECKER_TYPE_NEW_PICTURE         DuplicateCheckerType = 6
	DUPLICATE_CHECKER_TYPE_IMPORT_TRANSACTIONS DuplicateCheckerType = 7
	DUPLICATE_CHECKER_TYPE_FAILURE_CHECK       DuplicateCheckerType = 255
)
