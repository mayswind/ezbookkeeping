package extmw

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStampCommentInBody_AddsTheMarkAndKeepsEverythingElse(t *testing.T) {
	body := []byte(`{"type":3,"categoryId":"123456789012345678","time":1700000000,"sourceAmount":150000,"comment":"lunch","tagIds":["1","2"]}`)

	out, ok := StampCommentInBody(body, "Sam")
	require.True(t, ok)

	var fields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(out, &fields))

	assert.JSONEq(t, `"lunch · by Sam"`, string(fields["comment"]))
	assert.Equal(t, `"123456789012345678"`, string(fields["categoryId"]), "ids must not be turned into floats")
	assert.Equal(t, `150000`, string(fields["sourceAmount"]))
	assert.Equal(t, `["1","2"]`, string(fields["tagIds"]))
	assert.Equal(t, 6, len(fields))
}

func TestStampCommentInBody_NoCommentYetAndSpecialCharacters(t *testing.T) {
	out, ok := StampCommentInBody([]byte(`{"sourceAmount":5}`), `Zoë "Z" <b>`)
	require.True(t, ok)

	var parsed struct {
		Comment      string `json:"comment"`
		SourceAmount int    `json:"sourceAmount"`
	}
	require.NoError(t, json.Unmarshal(out, &parsed))
	assert.Equal(t, `by Zoë "Z" <b>`, parsed.Comment)
	assert.Equal(t, 5, parsed.SourceAmount)
}

func TestStampCommentInBody_RespectsTheCommentLengthLimit(t *testing.T) {
	long := make([]byte, 300)
	for i := range long {
		long[i] = 'x'
	}

	out, ok := StampCommentInBody([]byte(`{"comment":"`+string(long)+`"}`), "Sam")
	require.True(t, ok)

	var parsed struct{ Comment string }
	require.NoError(t, json.Unmarshal(out, &parsed))
	assert.LessOrEqual(t, len([]rune(parsed.Comment)), 255)
	assert.Contains(t, parsed.Comment, "by Sam")
}

func TestStampCommentInBody_LeavesOtherBodiesAlone(t *testing.T) {
	for _, body := range []string{``, `not json`, `[1,2]`, `null`, `{"comment":5}`} {
		out, ok := StampCommentInBody([]byte(body), "Sam")
		assert.False(t, ok, body)
		assert.Equal(t, body, string(out), body)
	}

	out, ok := StampCommentInBody([]byte(`{"comment":"x"}`), "")
	assert.True(t, ok)
	assert.JSONEq(t, `{"comment":"x"}`, string(out), "no name means no mark")
}

func TestDescribeRequest(t *testing.T) {
	cases := map[string][2]string{
		"/ext/sales/add.json":                     {"sales.add", "sales"},
		"/ext/stock/receive.json":                 {"stock.receive", "stock"},
		"/transactions/add.json":                  {"transactions.add", "transactions"},
		"/transactions/batch_update/account.json": {"transactions.batch_update.account", "batch_update"},
		"/transaction/categories/modify.json":     {"transaction.categories.modify", "categories"},
		"/accounts/hide.json":                     {"accounts.hide", "accounts"},
		"/ext/items/delete.json":                  {"items.delete", "items"},
	}

	for path, want := range cases {
		action, entity := DescribeRequest(path)
		assert.Equal(t, want[0], action, path)
		assert.Equal(t, want[1], entity, path)
	}

	action, entity := DescribeRequest("")
	assert.Equal(t, "", action)
	assert.Equal(t, "", entity)
}

func TestExtractEntityId(t *testing.T) {
	assert.Equal(t, int64(45), ExtractEntityId(nil, []byte(`{"result":{"id":"45","total":100},"success":true}`)), "created: from the response")
	assert.Equal(t, int64(12), ExtractEntityId([]byte(`{"id":"12"}`), []byte(`{"result":true,"success":true}`)), "changed: from the request")
	assert.Equal(t, int64(45), ExtractEntityId([]byte(`{"id":"12"}`), []byte(`{"result":{"id":"45"},"success":true}`)), "the response wins")
	assert.Equal(t, int64(7), ExtractEntityId([]byte(`{"id":7}`), nil), "plain numbers work too")
	assert.Equal(t, int64(0), ExtractEntityId([]byte(`{"name":"x"}`), []byte(`{"result":true}`)))
	assert.Equal(t, int64(31), ExtractEntityId(nil, []byte(`{"result":{"comment":"x","id":"31","uid":9},"success":true}`)), "id need not be the first field")
	assert.Equal(t, int64(0), ExtractEntityId(nil, []byte(`{"result":{"account":{"id":"5"}},"success":true}`)), "ids of nested things are not the entity's id")
	assert.Equal(t, int64(0), ExtractEntityId(nil, nil))
	assert.Equal(t, int64(0), ExtractEntityId([]byte(`{"id":"99999999999999999999999"}`), nil), "too large to be an id")
}
