package extservices

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

// AuditService stores and reads the audit log of delegated write requests
type AuditService struct{}

// Audit is the audit service singleton
var Audit = &AuditService{}

// Record stores an audit entry
func (s *AuditService) Record(c core.Context, entry *extmodels.AuditLog) error {
	entry.CreatedUnix = nowUnix()
	_, err := globalDB().NewSession(c).Insert(entry)

	return err
}

// List returns the latest audit entries of a business, newest first
func (s *AuditService) List(c core.Context, ownerUid int64, limit int, beforeAuditId int64) ([]*extmodels.AuditLog, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	entries := make([]*extmodels.AuditLog, 0, limit)
	sess := globalDB().NewSession(c).Where("owner_uid=?", ownerUid)

	if beforeAuditId > 0 {
		sess = sess.And("audit_id<?", beforeAuditId)
	}

	err := sess.OrderBy("audit_id desc").Limit(limit).Find(&entries)

	return entries, err
}
