package internal

import "time"

// A Record is one model's usage inside one hour, for one tool.
type Record struct {
	Tool   string
	Model  string
	Hour   time.Time
	Input  int64
	Output int64
	Cache  int64
}

func (r Record) Total() int64 { return r.Input + r.Output + r.Cache }

type Totals struct {
	Input  int64   `json:"input"`
	Output int64   `json:"output"`
	Cache  int64   `json:"cache"`
	Tokens int64   `json:"tokens"`
	USD    float64 `json:"usd"`
	// Priced is false when a model in the bucket has no published rate. The
	// usd field is then a lower bound, not a total.
	Priced bool `json:"priced"`
}

func (t *Totals) add(r Record) {
	t.Input += r.Input
	t.Output += r.Output
	t.Cache += r.Cache
	t.Tokens += r.Total()
	if cost, ok := estimate(r.Model, r.Input, r.Output, r.Cache); ok {
		t.USD += cost
		return
	}
	t.Priced = false
}

type ToolTotal struct {
	Tool   string `json:"tool"`
	Totals Totals `json:"totals"`
}

type ModelTotal struct {
	Model  string `json:"model"`
	Totals Totals `json:"totals"`
}

type DayTotal struct {
	Day    string `json:"day"`
	Totals Totals `json:"totals"`
}

// A Summary is the shape a page asks for.
type Summary struct {
	SyncedAt        string       `json:"synced_at,omitempty"`
	NextSyncAt      string       `json:"next_sync_at,omitempty"`
	IntervalSeconds int          `json:"interval_seconds,omitempty"`
	From            string       `json:"from"`
	To              string       `json:"to"`
	Totals          Totals       `json:"totals"`
	Tools           []ToolTotal  `json:"tools"`
	Models          []ModelTotal `json:"models"`
	Days            []DayTotal   `json:"days"`
}

func Summarize(records []Record, from, to time.Time) Summary {
	byTool := map[string]*Totals{}
	byModel := map[string]*Totals{}
	byDay := map[string]*Totals{}

	for _, r := range records {
		if r.Hour.Before(from) || !r.Hour.Before(to) {
			continue
		}
		add := func(m map[string]*Totals, key string) {
			t := m[key]
			if t == nil {
				t = &Totals{Priced: true}
				m[key] = t
			}
			t.add(r)
		}
		add(byTool, r.Tool)
		add(byModel, r.Model)
		add(byDay, r.Hour.Format("2006-01-02"))
	}

	out := Summary{From: from.Format(time.RFC3339), To: to.Format(time.RFC3339), Totals: Totals{Priced: true}}
	for tool, t := range byTool {
		out.Tools = append(out.Tools, ToolTotal{Tool: tool, Totals: *t})
		out.Totals.Input += t.Input
		out.Totals.Output += t.Output
		out.Totals.Cache += t.Cache
		out.Totals.Tokens += t.Tokens
		out.Totals.USD += t.USD
		if !t.Priced {
			out.Totals.Priced = false
		}
	}
	for model, t := range byModel {
		out.Models = append(out.Models, ModelTotal{Model: model, Totals: *t})
	}
	for day, t := range byDay {
		out.Days = append(out.Days, DayTotal{Day: day, Totals: *t})
	}
	sortSummary(&out)
	return out
}
