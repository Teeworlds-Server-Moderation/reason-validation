package main

import (
	"sort"
	"sync"

	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/jszwec/csvutil"
)

// CSet is a concurrent set that is used as key value store
type CSet struct {
	m  map[string]map[string]string
	mu sync.Mutex
}

// NewCSet creates a new Concurrent Set
func NewCSet() *CSet {
	return &CSet{
		m: make(map[string]map[string]string, 4096),
	}
}

// Size returns the current map size
func (cs *CSet) Size() int {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return len(cs.m)
}

// Add reason "fv" action "kickvote", reaction "voteban"
func (cs *CSet) Add(reason, action, reaction string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if cs.m[reason] == nil {
		cs.m[reason] = make(map[string]string, 2)
	}
	cs.m[reason][action] = reaction
}

// Get a reaction based on the reason and the action (kick or specvote)
func (cs *CSet) Get(reason, action string) (reaction string, ok bool) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	reactionMap, ok := cs.m[reason]
	if !ok {
		return "", false
	}

	reaction, ok = reactionMap[action]
	if !ok {
		return "", false
	}
	return reaction, true
}

func (cs *CSet) AddFromCSV(reason Reason) {
	cs.Add(reason.Reason, events.TypeVoteKickStarted, reason.KickvoteAction)
	cs.Add(reason.Reason, events.TypeVoteSpecStarted, reason.SpecVoteAction)
}

func (cs *CSet) DumpCSV() ([]byte, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	list := make([]Reason, 0, len(cs.m))

	for reason, reactionMap := range cs.m {
		specvoteReaction, ok := reactionMap[events.TypeVoteSpecStarted]
		if !ok {
			continue
		}
		kickvoteReaction, ok := reactionMap[events.TypeVoteSpecStarted]
		if !ok {
			continue
		}
		list = append(list, Reason{
			Reason:         reason,
			SpecVoteAction: specvoteReaction,
			KickvoteAction: kickvoteReaction,
		})
	}

	// sort before dumping
	sort.Sort(unknownAtFrontAndSortedByReason(list))

	return csvutil.Marshal(list)
}
