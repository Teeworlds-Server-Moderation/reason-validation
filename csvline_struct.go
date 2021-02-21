package main

// Reason represents one line in the csv file
type Reason struct {
	Reason         string `csv:"reason"`
	SpecVoteAction string `csv:"vote_spec_action"`
	KickvoteAction string `csv:"vote_kick_action"`
}
