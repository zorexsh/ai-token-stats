package internal

import "slices"

func sortSummary(s *Summary) {
	slices.SortFunc(s.Tools, func(a, b ToolTotal) int {
		return cmpDesc(a.Totals.Tokens, b.Totals.Tokens, a.Tool, b.Tool)
	})
	slices.SortFunc(s.Models, func(a, b ModelTotal) int {
		return cmpDesc(a.Totals.Tokens, b.Totals.Tokens, a.Model, b.Model)
	})
	slices.SortFunc(s.Days, func(a, b DayTotal) int {
		if a.Day < b.Day {
			return -1
		}
		if a.Day > b.Day {
			return 1
		}
		return 0
	})
}

func cmpDesc(av, bv int64, an, bn string) int {
	if av != bv {
		if av > bv {
			return -1
		}
		return 1
	}
	if an < bn {
		return -1
	}
	if an > bn {
		return 1
	}
	return 0
}
