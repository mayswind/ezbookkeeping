package models

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
)

func TestUserExternalAuthInfoResponsesSliceLess(t *testing.T) {
	var userExternalAuthInfoResponsesSlice UserExternalAuthInfoResponsesSlice
	userExternalAuthInfoResponsesSlice = append(userExternalAuthInfoResponsesSlice, &UserExternalAuthInfoResponse{
		ExternalAuthType: core.USER_EXTERNAL_AUTH_TYPE_OAUTH2_OIDC,
		Linked:           true,
		ExternalUsername: "test",
		CreatedAt:        int64(1),
	})
	userExternalAuthInfoResponsesSlice = append(userExternalAuthInfoResponsesSlice, &UserExternalAuthInfoResponse{
		ExternalAuthType: core.USER_EXTERNAL_AUTH_TYPE_OAUTH2_NEXTCLOUD,
		Linked:           false,
	})
	userExternalAuthInfoResponsesSlice = append(userExternalAuthInfoResponsesSlice, &UserExternalAuthInfoResponse{
		ExternalAuthType: core.USER_EXTERNAL_AUTH_TYPE_OAUTH2_GITEA,
		Linked:           false,
	})
	userExternalAuthInfoResponsesSlice = append(userExternalAuthInfoResponsesSlice, &UserExternalAuthInfoResponse{
		ExternalAuthType: core.USER_EXTERNAL_AUTH_TYPE_OAUTH2_GITHUB,
		Linked:           true,
		ExternalUsername: "test4",
		CreatedAt:        int64(2),
	})

	sort.Sort(userExternalAuthInfoResponsesSlice)

	assert.Equal(t, core.USER_EXTERNAL_AUTH_TYPE_OAUTH2_GITHUB, userExternalAuthInfoResponsesSlice[0].ExternalAuthType)
	assert.Equal(t, core.USER_EXTERNAL_AUTH_TYPE_OAUTH2_OIDC, userExternalAuthInfoResponsesSlice[1].ExternalAuthType)
	assert.Equal(t, core.USER_EXTERNAL_AUTH_TYPE_OAUTH2_GITEA, userExternalAuthInfoResponsesSlice[2].ExternalAuthType)
	assert.Equal(t, core.USER_EXTERNAL_AUTH_TYPE_OAUTH2_NEXTCLOUD, userExternalAuthInfoResponsesSlice[3].ExternalAuthType)
}
