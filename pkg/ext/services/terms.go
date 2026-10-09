package extservices

import (
	"regexp"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

// TermsService records who accepted which version of the Terms of Service and Privacy Policy.
// The current version is published with the app (legal/details.json); the server stores what the person accepted.
type TermsService struct{}

// Terms is the terms service singleton
var Terms = &TermsService{}

var versionPattern = regexp.MustCompile(`^[A-Za-z0-9._-]{1,32}$`)

// Latest returns the version a person accepted most recently and when, or "" if they never did
func (s *TermsService) Latest(c core.Context, uid int64) (version string, acceptedUnix int64, err error) {
	if uid <= 0 {
		return "", 0, errs.ErrUserIdInvalid
	}

	row := &extmodels.TermsAcceptance{}
	has, err := globalDB().NewSession(c).Where("uid=?", uid).OrderBy("acceptance_id desc").Get(row)

	if err != nil || !has {
		return "", 0, err
	}

	return row.Version, row.AcceptedUnix, nil
}

// Accept records that a person accepted a version. Accepting the version they already accepted last changes nothing,
// so a double click or a retry does not add rows.
func (s *TermsService) Accept(c core.Context, uid int64, version string, clientIp string) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if !versionPattern.MatchString(version) {
		return exterrs.ErrTermsVersionInvalid
	}

	if len(clientIp) > 45 {
		clientIp = clientIp[:45]
	}

	return globalDB().DoTransaction(c, func(sess *xormSession) error {
		latest := &extmodels.TermsAcceptance{}
		has, err := sess.Where("uid=?", uid).OrderBy("acceptance_id desc").Get(latest)

		if err != nil {
			return err
		} else if has && latest.Version == version {
			return nil
		}

		_, err = sess.Insert(&extmodels.TermsAcceptance{Uid: uid, Version: version, AcceptedUnix: nowUnix(), ClientIp: clientIp})

		return err
	})
}
