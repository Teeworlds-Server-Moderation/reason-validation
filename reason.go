package main

// Reason represents one line in the csv file
type Reason struct {
	Reason         string `csv:"reason"`
	SpecVoteAction string `csv:"vote_spec_action"`
	KickvoteAction string `csv:"vote_kick_action"`
}

type unknownAtFrontAndSortedByReason []Reason

func (a unknownAtFrontAndSortedByReason) Len() int      { return len(a) }
func (a unknownAtFrontAndSortedByReason) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a unknownAtFrontAndSortedByReason) Less(i, j int) bool {
	// put unknown actions
	if a[i].KickvoteAction == "unknown" && a[j].KickvoteAction != "unknown" {
		return true
	} else if a[i].KickvoteAction != "unknown" && a[j].KickvoteAction == "unknown" {
		return false
	}

	return a[i].Reason < a[j].Reason
}
