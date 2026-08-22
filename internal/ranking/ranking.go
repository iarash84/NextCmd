package ranking

import (
	"nextcmd/sdk"
	"sort"
	"strings"
)

func Rank(input string, suggestions []sdk.Suggestion, limit int) []sdk.Suggestion {
	type scored struct {
		suggestion sdk.Suggestion
		score      int
	}
	unique := map[string]scored{}
	query := strings.ToLower(strings.TrimSpace(input))
	for _, suggestion := range suggestions {
		display := strings.ToLower(suggestion.Command.Display())
		score := suggestion.Priority*10 + suggestion.Relevance
		if query == "" || strings.HasPrefix(display, query) {
			score += 1000
		} else if fuzzy(query, display) {
			score += 300
		} else {
			continue
		}
		key := suggestion.Command.Display()
		candidate := scored{suggestion, score}
		if old, ok := unique[key]; !ok || candidate.score > old.score {
			unique[key] = candidate
		}
	}
	items := make([]scored, 0, len(unique))
	for _, item := range unique {
		items = append(items, item)
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].score != items[j].score {
			return items[i].score > items[j].score
		}
		return items[i].suggestion.Command.Display() < items[j].suggestion.Command.Display()
	})
	if limit <= 0 || limit > len(items) {
		limit = len(items)
	}
	out := make([]sdk.Suggestion, limit)
	for i := range out {
		out[i] = items[i].suggestion
	}
	return out
}
func fuzzy(needle, haystack string) bool {
	i := 0
	for _, r := range haystack {
		if i < len(needle) && byte(r) == needle[i] {
			i++
		}
	}
	return i == len(needle)
}
