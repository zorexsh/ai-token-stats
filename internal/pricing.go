package internal

import "strings"

// A rate is USD per million tokens. Cache reads bill at a tenth of input and
// cache writes at 1.25x where the provider publishes both, which is the same
// split Token Tracker uses. Rates were checked against provider pages in
// November 2026. An unknown model is left unpriced rather than guessed.
type rate struct {
	input  float64
	output float64
	cache  float64
	write  float64
}

var rates = map[string]rate{
	"claude-opus-5-5":   {4, 20, 0.2, 5},
	"claude-opus-5":     {5, 25, 0.5, 6.25},
	"claude-opus-4-8":   {5, 25, 0.5, 6.25},
	"claude-opus-4-7":   {5, 25, 0.5, 6.25},
	"claude-opus-4-6":   {5, 25, 0.5, 6.25},
	"claude-opus-4-5":   {5, 25, 0.5, 6.25},
	"claude-sonnet-5-5": {2, 10, 0.1, 2.5},
	"claude-sonnet-5":   {2, 10, 0.2, 2.5},
	"claude-sonnet-4-6": {3, 15, 0.3, 3.75},
	"claude-sonnet-4-5": {3, 15, 0.3, 3.75},
	"claude-haiku-5-5":  {0.1, 0.5, 0.01, 0.125},
	"claude-haiku-4-5":  {1, 5, 0.1, 1.25},
	"claude-fable-5":    {10, 50, 1, 12.5},

	"gpt-5.6":       {4, 20, 0.4, 5},
	"gpt-5.6-sol":   {4, 20, 0.4, 5},
	"gpt-5.6-terra": {2, 12, 0.2, 2.5},
	"gpt-5.6-luna":  {0.2, 1.2, 0.02, 0.25},
	"gpt-5.5":       {5, 30, 0.5, 0},
	"gpt-5.4":       {2.5, 15, 0.25, 0},
	"gpt-5.4-pro":   {30, 180, 3, 0},
	"gpt-5.4-mini":  {0.75, 4.5, 0.075, 0},
	"gpt-5":         {1.25, 10, 0.125, 0},
	"gpt-5-mini":    {0.25, 2, 0.025, 0},
	"gpt-4o":        {2.5, 10, 1.25, 0},
	"gpt-4o-mini":   {0.15, 0.6, 0.075, 0},

	"gemini-2.5-pro":   {1.25, 10, 0.125, 0},
	"gemini-2.5-flash": {0.3, 2.5, 0.03, 0},
	"gemini-3-pro":     {2, 12, 0.2, 0},
	"gemini-4-argon":   {2, 10, 0.1, 0},

	"deepseek-v4-pro":   {1.32, 3.96, 0.044, 1.32},
	"deepseek-v4-flash": {0.44, 1.32, 0.014, 0.44},
	"deepseek-chat":     {0.14, 0.28, 0.0028, 0.14},
	"deepseek-reasoner": {0.14, 0.28, 0.0028, 0.14},

	"kimi-k3":   {3, 15, 0.3, 0},
	"kimi-k2.7": {0.95, 4, 0.19, 0},
	"kimi-k2.6": {0.95, 4, 0.16, 0},
	"kimi-k2.5": {0.6, 2, 0.15, 0},

	"glm-5.3":       {1.4, 4.4, 0.26, 0},
	"glm-5.3-flash": {0.15, 0.5, 0.03, 0},
	"glm-5.2":       {1.4, 4.4, 0.26, 0},
	"glm-5.1":       {1.4, 4.4, 0.26, 0},
	"glm-5":         {1, 3.2, 0.2, 0},
	"glm-4.7":       {0.6, 2.2, 0.11, 0},
	"glm-4.6":       {0.6, 2.2, 0.11, 0},
	"glm-4.5":       {0.6, 2.2, 0.11, 0},

	"grok-4":   {3, 15, 0.75, 0},
	"grok-4.5": {2, 6, 0.5, 0},
	"grok-4.6": {2, 6, 0.5, 0},
	"grok-4.7": {2, 6, 0.5, 0},

	"composer-2":   {0.5, 2.5, 0.2, 0},
	"composer-1.5": {3.5, 17.5, 0.35, 0},
	"composer-1":   {1.25, 10, 0.125, 0},
}

func init() {
	normalized := make(map[string]rate, len(rates))
	for key, r := range rates {
		normalized[strings.ReplaceAll(key, ".", "-")] = r
	}
	rates = normalized
}

// freeMarkers are substrings that mean the route bills nothing per token.
// A gateway's "free" suffix is a billing fact, not a guess. A fine-tune or
// abliteration suffix is not, so those fall through to the base model.
var freeMarkers = []string{"-free", ":free"}

// estimate returns the USD cost and whether the model was recognized. An
// unrecognized model returns false so the caller can omit the number instead
// of inventing one.
func estimate(model string, input, output, cache int64) (float64, bool) {
	r, ok := lookup(model)
	if !ok {
		return 0, false
	}
	cost := (float64(input)*r.input + float64(output)*r.output + float64(cache)*r.cache) / 1_000_000
	return cost, true
}

func lookup(model string) (rate, bool) {
	name := normalize(model)
	if name == "" {
		return rate{}, false
	}
	if isFree(name) {
		return rate{}, true
	}
	if r, ok := rates[name]; ok {
		return r, true
	}
	// Gateways append a variant the price card does not list. Take the
	// longest known prefix, so "glm-5-3-flash-abliterated" prices as the
	// flash tier rather than the base model.
	var best string
	var found rate
	for key, r := range rates {
		if strings.HasPrefix(name, key+"-") && len(key) > len(best) {
			best = key
			found = r
		}
	}
	if best == "" {
		return rate{}, false
	}
	return found, true
}

func normalize(model string) string {
	name := strings.ToLower(strings.TrimSpace(model))
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	name = strings.TrimPrefix(name, "fw-")
	return strings.ReplaceAll(name, ".", "-")
}

func isFree(name string) bool {
	for _, marker := range freeMarkers {
		if strings.Contains(name, marker) {
			return true
		}
	}
	return false
}
