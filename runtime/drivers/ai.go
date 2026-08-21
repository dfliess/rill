package drivers

import (
	"context"

	"github.com/google/jsonschema-go/jsonschema"
	aiv1 "github.com/rilldata/rill/proto/gen/rill/ai/v1"
)

type AIService interface {
	Complete(ctx context.Context, opts *CompleteOptions) (*CompleteResult, error)
}

type CompleteOptions struct {
	Messages     []*aiv1.CompletionMessage
	Tools        []*aiv1.Tool
	OutputSchema *jsonschema.Schema
}

type CompleteResult struct {
	Message *aiv1.CompletionMessage
	// Provider is the LLM provider that served the completion (e.g. "claude", "openai", "gemini"). For the managed admin
	// proxy this is the real underlying provider, not "admin", so it stays meaningful through the proxy.
	Provider string
	// Model is the specific model that served the completion (e.g. "gemini-3.5-flash-lite"). Provider alone cannot price
	// a completion: models of one provider differ in rate by an order of magnitude, so a deployment that runs a large
	// model for some agents and a small one for others cannot be billed from the provider tag. Empty when the driver
	// cannot report it, so consumers must treat it as optional rather than assume a value.
	Model string
	// InputTokens and CachedInputTokens are reported as the provider returns them, so their relationship is
	// provider-specific: for Claude they are disjoint (InputTokens excludes cached); for OpenAI/Gemini CachedInputTokens
	// is a subset of InputTokens.
	InputTokens       int
	CachedInputTokens int
	OutputTokens      int
}
