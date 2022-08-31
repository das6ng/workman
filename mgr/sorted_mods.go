package mgr

import "sort"

type sortedMods struct {
	used   map[string]struct{}
	merged []string
}

func loadSortedMods(used, nearby []string) sortedMods {
	m := make(map[string]struct{}, len(used))
	for _, u := range used {
		m[u] = struct{}{}
	}
	return sortedMods{
		used:   m,
		merged: uniqStrList(append(used, nearby...)),
	}
}

func (x sortedMods) Mods() []string {
	sort.Sort(x)
	return x.merged
}

func (x sortedMods) Len() int {
	return len(x.merged)
}
func (x sortedMods) Less(i, j int) bool {
	_, ui := x.used[x.merged[i]]
	_, uj := x.used[x.merged[j]]
	return (ui && !uj) || (ui && uj || !ui && !uj) && x.merged[i] < x.merged[j]
}
func (x sortedMods) Swap(i, j int) {
	x.merged[i], x.merged[j] = x.merged[j], x.merged[i]
}
