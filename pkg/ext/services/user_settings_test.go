package extservices

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

func TestUserSettings_DefaultIsOffAndNotConfigured(t *testing.T) {
	f := newFixture(t)

	setting, configured, err := UserSettings.Get(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.False(t, configured)
	assert.False(t, setting.BusinessFeatures, "features are opt-in")
	assert.Equal(t, f.ownerUid, setting.Uid)
}

func TestUserSettings_TurnOnAndOffPersists(t *testing.T) {
	f := newFixture(t)

	on, err := UserSettings.SetBusinessFeatures(f.c, f.ownerUid, true)
	require.NoError(t, err)
	assert.True(t, on.BusinessFeatures)

	got, configured, err := UserSettings.Get(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.True(t, configured)
	assert.True(t, got.BusinessFeatures)

	off, err := UserSettings.SetBusinessFeatures(f.c, f.ownerUid, false)
	require.NoError(t, err)
	assert.False(t, off.BusinessFeatures)

	got, configured, err = UserSettings.Get(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.True(t, configured, "an explicit 'off' is a choice too, not the same as never having chosen")
	assert.False(t, got.BusinessFeatures)
}

func TestUserSettings_AreSetOnceNotDuplicated(t *testing.T) {
	f := newFixture(t)

	for i := 0; i < 3; i++ {
		_, err := UserSettings.SetBusinessFeatures(f.c, f.ownerUid, i%2 == 0)
		require.NoError(t, err)
	}

	rows := make([]*extmodels.UserSetting, 0)
	require.NoError(t, globalDB().NewSession(f.c).Where("uid=?", f.ownerUid).Find(&rows))
	assert.Len(t, rows, 1, "one row per person, however often the setting changes")

	got, _, err := UserSettings.Get(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.True(t, got.BusinessFeatures, "the last write (true) wins")
}

func TestUserSettings_PeopleAreIndependent(t *testing.T) {
	f := newFixture(t)
	other := f.newUser(t, "other")

	_, err := UserSettings.SetBusinessFeatures(f.c, f.ownerUid, true)
	require.NoError(t, err)

	got, configured, err := UserSettings.Get(f.c, other)
	require.NoError(t, err)
	assert.False(t, configured)
	assert.False(t, got.BusinessFeatures, "somebody else's choice must not leak")
}

func TestUserSettings_RejectInvalidUser(t *testing.T) {
	f := newFixture(t)

	_, _, err := UserSettings.Get(f.c, 0)
	assert.Equal(t, errs.ErrUserIdInvalid, err)
	_, err = UserSettings.SetBusinessFeatures(f.c, -1, true)
	assert.Equal(t, errs.ErrUserIdInvalid, err)
}
