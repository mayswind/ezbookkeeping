package duplicatechecker

// DuplicateChecker is common duplicate checker interface
type DuplicateChecker interface {
	GetSubmissionRemark(checkerType DuplicateCheckerType, uid int64, identification string) (bool, string)
	SetSubmissionRemark(checkerType DuplicateCheckerType, uid int64, identification string, remark string)
}
