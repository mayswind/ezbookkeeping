package extmw

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	extservices "github.com/mayswind/ezbookkeeping/pkg/ext/services"
)

const (
	maxCapturedRequestBody  = 1 << 20 // larger bodies (file uploads) are left alone
	maxCapturedResponseHead = 16384
	commentLimit            = 255
)

var (
	responseIdPattern = regexp.MustCompile(`"result"\s*:\s*\{\s*"id"\s*:\s*"?(\d+)"?`)
	requestIdPattern  = regexp.MustCompile(`"id"\s*:\s*"?(\d+)"?`)
)

// StampCommentInBody adds "by <name>" to the "comment" of a JSON request body, leaving every other field exactly
// as it was sent. It returns the body unchanged and false when the body is not a JSON object.
func StampCommentInBody(body []byte, name string) ([]byte, bool) {
	var fields map[string]json.RawMessage

	if err := json.Unmarshal(body, &fields); err != nil || fields == nil {
		return body, false
	}

	var comment string

	if raw, exists := fields["comment"]; exists {
		if err := json.Unmarshal(raw, &comment); err != nil {
			return body, false
		}
	}

	stamped, err := json.Marshal(extservices.StampText(comment, name, commentLimit))

	if err != nil {
		return body, false
	}

	fields["comment"] = stamped
	out, err := json.Marshal(fields)

	if err != nil {
		return body, false
	}

	return out, true
}

// DescribeRequest turns a route such as "/ext/sales/add.json" into an action ("sales.add") and the kind of thing it
// works on ("sales"), for the audit log.
func DescribeRequest(relativePath string) (action string, entityType string) {
	segments := strings.Split(strings.Trim(strings.TrimSuffix(relativePath, ".json"), "/"), "/")

	if len(segments) > 0 && segments[0] == "ext" {
		segments = segments[1:]
	}

	if len(segments) == 0 || segments[0] == "" {
		return "", ""
	}

	entityType = segments[0]

	if len(segments) >= 2 {
		entityType = segments[len(segments)-2]
	}

	action = strings.Join(segments, ".")

	if len(action) > 64 {
		action = action[:64]
	}

	if len(entityType) > 32 {
		entityType = entityType[:32]
	}

	return action, entityType
}

// ExtractEntityId finds the id of what a request created (in the response) or changed (in the request), or 0 when
// there is none.
func ExtractEntityId(requestBody []byte, responseHead []byte) int64 {
	if id := idFromResponse(responseHead); id > 0 {
		return id
	}

	if match := requestIdPattern.FindSubmatch(requestBody); len(match) == 2 {
		return parseId(match[1])
	}

	return 0
}

// idFromResponse reads result.id from a complete JSON response. A response too long to have been captured whole
// falls back to looking for "id" as the first field of the result.
func idFromResponse(head []byte) int64 {
	var parsed struct {
		Result json.RawMessage `json:"result"`
	}

	if err := json.Unmarshal(head, &parsed); err == nil {
		var result map[string]json.RawMessage

		if json.Unmarshal(parsed.Result, &result) == nil {
			return parseId(bytes.Trim(result["id"], `"`))
		}

		return 0
	}

	if match := responseIdPattern.FindSubmatch(head); len(match) == 2 {
		return parseId(match[1])
	}

	return 0
}

func parseId(text []byte) int64 {
	id, err := strconv.ParseInt(string(text), 10, 64)

	if err != nil || id < 0 {
		return 0
	}

	return id
}

// captureWriter remembers the start of a response while passing it through untouched
type captureWriter struct {
	gin.ResponseWriter
	head bytes.Buffer
}

func (w *captureWriter) keep(data []byte) {
	if room := maxCapturedResponseHead - w.head.Len(); room > 0 {
		if len(data) > room {
			data = data[:room]
		}

		w.head.Write(data)
	}
}

func (w *captureWriter) Write(data []byte) (int, error) {
	w.keep(data)
	return w.ResponseWriter.Write(data)
}

func (w *captureWriter) WriteString(data string) (int, error) {
	w.keep([]byte(data))
	return w.ResponseWriter.WriteString(data)
}
