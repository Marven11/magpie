package tune

import "testing"

// An agent's list names a model by its provider and its own name — Codex's
// "p1/claude-opus-5", "magpie/p1/claude-opus-5", a group's
// "group/auto-claude-opus-5" — and the vendor is the model's name, not the
// string's first letters. Read by prefix, every Claude model an agent reached
// through magpie was weighed as OpenAI's, so the window the card suggests and
// the tokens it says a change saves were both from another vendor's prices.
func TestWeightsForReadsTheModelsOwnName(t *testing.T) {
	anthropic := Weights{Read: 0.1, Write: 1.25, Write1h: 2, Out: 5}
	openai := Weights{Read: 0.1, Write: 1, Write1h: 1, Out: 8}
	for model, want := range map[string]Weights{
		"claude-opus-5":               anthropic,
		"p1/claude-opus-5":            anthropic,
		"magpie/p1/claude-opus-5":     anthropic,
		"group/auto-claude-opus-5":    anthropic,
		"CLAUDE-SONNET-4":             anthropic,
		"p1/Claude-Haiku-4":           anthropic,
		"gpt-5-codex":                 openai,
		"p1/gpt-5-codex":              openai,
		"magpie/p1/gpt-5-codex":       openai,
		"openrouter/anthropic/gemini": openai,
		"":                            openai,
	} {
		if got := WeightsFor(model); got != want {
			t.Errorf("WeightsFor(%q) = %+v, want %+v", model, got, want)
		}
	}
}
