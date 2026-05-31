package modelsdev

// Provider is a single inference provider listed on models.dev (OpenAI,
// Anthropic, Google, …) together with the models it serves.
type Provider struct {
	// ID is the models.dev provider identifier, e.g. "openai", "amazon-bedrock".
	ID string `json:"id"`
	// Name is the human-readable provider name, e.g. "OpenAI".
	Name string `json:"name"`
	// NPM is the @ai-sdk package for the provider, when one exists.
	NPM string `json:"npm,omitempty"`
	// Doc is a link to the provider's model documentation.
	Doc string `json:"doc,omitempty"`
	// Env lists the environment variables the provider's SDK reads for credentials.
	Env []string `json:"env,omitempty"`
	// Models is keyed by model ID (the same value as Model.ID).
	Models map[string]Model `json:"models"`
}

// Model is a single model offered by a Provider.
type Model struct {
	// ID is the provider-local model identifier, e.g. "gpt-4o".
	ID string `json:"id"`
	// Name is the human-readable model name, e.g. "GPT-4o".
	Name string `json:"name"`
	// Family groups related models, e.g. "claude-opus".
	Family string `json:"family,omitempty"`
	// Provider is the models.dev provider ID this model belongs to.
	Provider string `json:"provider,omitempty"`

	// Attachment reports whether the model accepts file/image attachments.
	Attachment bool `json:"attachment"`
	// Reasoning reports whether the model supports extended thinking / reasoning.
	Reasoning bool `json:"reasoning"`
	// ToolCall reports whether the model supports function/tool calling.
	ToolCall bool `json:"tool_call"`
	// Temperature reports whether the model honors a temperature parameter.
	Temperature bool `json:"temperature"`
	// OpenWeights reports whether the model's weights are openly available.
	OpenWeights bool `json:"open_weights"`
	// Interleaved reports whether the model supports interleaved thinking.
	Interleaved bool `json:"interleaved,omitempty"`
	// Experimental flags a model as experimental on models.dev.
	Experimental bool `json:"experimental,omitempty"`

	// StructuredOutput reports native JSON-schema structured output support.
	// It is a pointer because models.dev omits it for models where it is unknown.
	StructuredOutput *bool `json:"structured_output,omitempty"`

	// Status is the lifecycle status reported upstream (e.g. "deprecated"), if any.
	Status string `json:"status,omitempty"`
	// Knowledge is the training-data cutoff, e.g. "2024-10".
	Knowledge string `json:"knowledge,omitempty"`
	// ReleaseDate is the model's release date, e.g. "2026-05-28".
	ReleaseDate string `json:"release_date,omitempty"`
	// LastUpdated is when the model entry was last updated upstream.
	LastUpdated string `json:"last_updated,omitempty"`

	// Modalities lists the input/output modalities the model supports.
	Modalities Modalities `json:"modalities"`
	// Limit holds the context/output token limits.
	Limit Limit `json:"limit"`
	// Cost holds the per-million-token pricing in USD.
	Cost Cost `json:"cost"`
}

// Modalities describes the input and output modalities of a model. Values are
// lowercase tokens such as "text", "image", "audio", "pdf", "video".
type Modalities struct {
	Input  []string `json:"input,omitempty"`
	Output []string `json:"output,omitempty"`
}

// Limit holds a model's token limits. Zero means the limit is not reported.
type Limit struct {
	// Context is the maximum combined input+output context window in tokens.
	Context int `json:"context,omitempty"`
	// Input is the maximum input tokens, when reported separately.
	Input int `json:"input,omitempty"`
	// Output is the maximum output tokens.
	Output int `json:"output,omitempty"`
}

// Cost holds a model's pricing in USD per 1M tokens. All fields are pointers
// because models.dev omits the ones a provider does not publish; a nil field
// means "not reported", which is distinct from a reported zero price.
type Cost struct {
	// Input is the price per 1M input (prompt) tokens.
	Input *float64 `json:"input,omitempty"`
	// Output is the price per 1M output (completion) tokens.
	Output *float64 `json:"output,omitempty"`
	// CacheRead is the price per 1M prompt tokens served from cache (cache hits).
	CacheRead *float64 `json:"cache_read,omitempty"`
	// CacheWrite is the price per 1M tokens written to the prompt cache.
	CacheWrite *float64 `json:"cache_write,omitempty"`
	// Reasoning is the price per 1M reasoning tokens, when billed separately.
	Reasoning *float64 `json:"reasoning,omitempty"`
	// InputAudio / OutputAudio are per-1M-token audio prices, when applicable.
	InputAudio  *float64 `json:"input_audio,omitempty"`
	OutputAudio *float64 `json:"output_audio,omitempty"`
	// ContextOver200K is the input price for context beyond 200K tokens, when
	// the provider applies long-context pricing.
	ContextOver200K *float64 `json:"context_over_200k,omitempty"`
}

// SupportsInput reports whether the model accepts the given input modality.
func (m Model) SupportsInput(modality string) bool {
	for _, in := range m.Modalities.Input {
		if in == modality {
			return true
		}
	}
	return false
}

// SupportsVision reports whether the model accepts image input.
func (m Model) SupportsVision() bool { return m.SupportsInput("image") }
