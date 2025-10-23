package github

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

func TestGithubOAuth2Datasource_GetUserInfoRequest(t *testing.T) {
	datasource := &GithubOAuth2DataSource{}
	req, err := datasource.GetUserInfoRequest()

	assert.Nil(t, err)
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "https://api.github.com/user", req.URL.String())
	assert.Equal(t, "application/vnd.github+json", req.Header.Get("Accept"))
}

func TestGithubOAuth2Datasource_ParseUserInfo_Success(t *testing.T) {
	datasource := &GithubOAuth2DataSource{}
	responseContent := `{
		"login": "octocat",
		"id": 1,
		"node_id": "MDQ6VXNlcjE=",
		"avatar_url": "https://github.com/images/error/octocat_happy.gif",
		"gravatar_id": "",
		"url": "https://api.github.com/users/octocat",
		"html_url": "https://github.com/octocat",
		"followers_url": "https://api.github.com/users/octocat/followers",
		"following_url": "https://api.github.com/users/octocat/following{/other_user}",
		"gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
		"starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
		"subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
		"organizations_url": "https://api.github.com/users/octocat/orgs",
		"repos_url": "https://api.github.com/users/octocat/repos",
		"events_url": "https://api.github.com/users/octocat/events{/privacy}",
		"received_events_url": "https://api.github.com/users/octocat/received_events",
		"type": "User",
		"site_admin": false,
		"name": "monalisa octocat",
		"company": "GitHub",
		"blog": "https://github.com/blog",
		"location": "San Francisco",
		"email": "octocat@github.com",
		"hireable": false,
		"bio": "There once was...",
		"twitter_username": "monatheoctocat",
		"public_repos": 2,
		"public_gists": 1,
		"followers": 20,
		"following": 0,
		"created_at": "2008-01-14T04:33:35Z",
		"updated_at": "2008-01-14T04:33:35Z",
		"private_gists": 81,
		"total_private_repos": 100,
		"owned_private_repos": 100,
		"disk_usage": 10000,
		"collaborators": 8,
		"two_factor_authentication": true,
		"plan": {
			"name": "Medium",
			"space": 400,
			"private_repos": 20,
			"collaborators": 0
		}
	}`
	info, err := datasource.ParseUserInfo(core.NewNullContext(), []byte(responseContent))

	assert.Nil(t, err)
	assert.Equal(t, "octocat", info.UserName)
	assert.Equal(t, "octocat@github.com", info.Email)
	assert.Equal(t, "monalisa octocat", info.NickName)
}

func TestGithubOAuth2Datasource_ParseUserInfo_InvalidJson(t *testing.T) {
	datasource := &GithubOAuth2DataSource{}
	_, err := datasource.ParseUserInfo(core.NewNullContext(), []byte("invalid"))

	assert.Equal(t, errs.ErrCannotRetrieveUserInfo, err)
}

func TestGithubOAuth2Datasource_ParseUserInfo_EmptyLogin(t *testing.T) {
	datasource := &GithubOAuth2DataSource{}
	responseContent := `{"login": ""}`
	_, err := datasource.ParseUserInfo(core.NewNullContext(), []byte(responseContent))

	assert.Equal(t, errs.ErrCannotRetrieveUserInfo, err)
}
