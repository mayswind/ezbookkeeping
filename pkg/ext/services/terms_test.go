package extservices

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

func (f *fixture) acceptanceRows(t *testing.T, uid int64) []*extmodels.TermsAcceptance {
	t.Helper()

	rows := make([]*extmodels.TermsAcceptance, 0)
	require.NoError(t, globalDB().NewSession(f.c).Where("uid=?", uid).OrderBy("acceptance_id").Find(&rows))

	return rows
}

func TestTerms_NobodyHasAcceptedAtFirst(t *testing.T) {
	f := newFixture(t)

	version, at, err := Terms.Latest(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.Equal(t, "", version)
	assert.Equal(t, int64(0), at)
}

func TestTerms_AcceptingIsRecordedWithTimeAndAddress(t *testing.T) {
	f := newFixture(t)
	original := nowUnix
	nowUnix = func() int64 { return 1_800_000_000 }
	t.Cleanup(func() { nowUnix = original })

	require.NoError(t, Terms.Accept(f.c, f.ownerUid, "2026-10-09", "203.0.113.5"))

	version, at, err := Terms.Latest(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.Equal(t, "2026-10-09", version)
	assert.Equal(t, int64(1_800_000_000), at)

	rows := f.acceptanceRows(t, f.ownerUid)
	require.Len(t, rows, 1)
	assert.Equal(t, "203.0.113.5", rows[0].ClientIp)
}

func TestTerms_AcceptingTheSameVersionTwiceAddsNothing(t *testing.T) {
	f := newFixture(t)

	for i := 0; i < 3; i++ {
		require.NoError(t, Terms.Accept(f.c, f.ownerUid, "v1", "198.51.100.1"))
	}

	assert.Len(t, f.acceptanceRows(t, f.ownerUid), 1, "a double click or a retry must not add rows")
}

func TestTerms_ANewVersionKeepsTheHistory(t *testing.T) {
	f := newFixture(t)

	require.NoError(t, Terms.Accept(f.c, f.ownerUid, "v1", ""))
	require.NoError(t, Terms.Accept(f.c, f.ownerUid, "v2", ""))

	version, _, err := Terms.Latest(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.Equal(t, "v2", version)

	rows := f.acceptanceRows(t, f.ownerUid)
	require.Len(t, rows, 2)
	assert.Equal(t, []string{"v1", "v2"}, []string{rows[0].Version, rows[1].Version}, "what was accepted before stays on record")

	// going back to an older version is a new acceptance, not a rewrite
	require.NoError(t, Terms.Accept(f.c, f.ownerUid, "v1", ""))
	assert.Len(t, f.acceptanceRows(t, f.ownerUid), 3)
}

func TestTerms_PeopleAreIndependent(t *testing.T) {
	f := newFixture(t)
	other := f.newUser(t, "other")

	require.NoError(t, Terms.Accept(f.c, f.ownerUid, "v1", ""))

	version, _, err := Terms.Latest(f.c, other)
	require.NoError(t, err)
	assert.Equal(t, "", version, "somebody else's acceptance counts for nobody else")
}

func TestTerms_BadInputIsRefused(t *testing.T) {
	f := newFixture(t)

	for _, bad := range []string{"", "has space", "<script>", strings.Repeat("v", 33), "quote'"} {
		assert.Equal(t, exterrs.ErrTermsVersionInvalid, Terms.Accept(f.c, f.ownerUid, bad, ""), bad)
	}

	assert.Equal(t, errs.ErrUserIdInvalid, Terms.Accept(f.c, 0, "v1", ""))
	assert.Empty(t, f.acceptanceRows(t, f.ownerUid))

	// a long address is cut to fit its column rather than failing the acceptance
	require.NoError(t, Terms.Accept(f.c, f.ownerUid, "v1", strings.Repeat("1", 80)))
	assert.Len(t, f.acceptanceRows(t, f.ownerUid)[0].ClientIp, 45)
}
