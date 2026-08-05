package openai

// OpenAIMessageRole defines the role of OpenAI chat completions message
type OpenAIMessageRole string

// OpenAI Message Roles
const (
	OpenAIMessageRoleSystem OpenAIMessageRole = "system"
	OpenAIMessageRoleUser   OpenAIMessageRole = "user"
)
