package duplicatechecker

// DuplicateChecker is common duplicate checker interface
type DuplicateChecker interface {
	Get(checkerType DuplicateCheckerType, uid int64, identification string) (bool, string)
	Set(checkerType DuplicateCheckerType, uid int64, identification string, remark string)
}
