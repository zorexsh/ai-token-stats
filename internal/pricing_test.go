package internal

import "testing"

func TestEstimateKnownModel(t *testing.T) {
	got, ok := estimate("claude-opus-4-6", 1_000_000, 0, 0)
	if !ok || got != 5 {
		t.Fatalf("cost = %v ok=%v", got, ok)
	}
	got, ok = estimate("7api-claude/claude-opus-4-6", 0, 1_000_000, 0)
	if !ok || got != 25 {
		t.Fatalf("provider prefix cost = %v ok=%v", got, ok)
	}
}

func TestEstimateCacheReadsCheaper(t *testing.T) {
	got, ok := estimate("claude-opus-4-6", 0, 0, 1_000_000)
	if !ok || got != 0.5 {
		t.Fatalf("cache cost = %v, want 0.5", got)
	}
}

func TestEstimateFreeRouteIsZero(t *testing.T) {
	got, ok := estimate("opencode/deepseek-v4-flash-free", 5_000_000, 1_000_000, 0)
	if !ok || got != 0 {
		t.Fatalf("free route cost = %v ok=%v", got, ok)
	}
}

func TestEstimateWalksBackToKnownPrefix(t *testing.T) {
	name := normalize("7api-other/glm-5.3-flash-abliterated")
	if name != "glm-5-3-flash-abliterated" {
		t.Fatalf("normalized = %q", name)
	}
	got, ok := estimate("7api-other/glm-5.3-flash-abliterated", 1_000_000, 0, 0)
	want := rates["glm-5-3-flash"].input
	if !ok || got != want {
		t.Fatalf("prefix cost = %v ok=%v, want %v", got, ok, want)
	}
}

func TestEstimateUnknownModel(t *testing.T) {
	if _, ok := estimate("ollama/qwen3.5:9b", 100, 100, 0); ok {
		t.Fatal("unknown model was priced")
	}
}
