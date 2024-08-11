package duplicatechecker

// DuplicateCheckerType represents duplicate checker type
type DuplicateCheckerType uint8

// Types of uuid
const (
	DUPLICATE_CHECKER_TYPE_BACKGROUND_CRON_JOB DuplicateCheckerType = 0
	DUPLICATE_CHECKER_TYPE_NEW_ACCOUNT         DuplicateCheckerType = 1
	DUPLICATE_CHECKER_TYPE_NEW_CATEGORY        DuplicateCheckerType = 2
	DUPLICATE_CHECKER_TYPE_NEW_TRANSACTION     DuplicateCheckerType = 3
	DUPLICATE_CHECKER_TYPE_NEW_TEMPLATE        DuplicateCheckerType = 4
)
