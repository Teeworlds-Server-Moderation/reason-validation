package main

import (
	"sync"

	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/jszwec/csvutil"
)

type CSet struct {
	m  map[string]map[string]string
	mu sync.Mutex
}

func NewCSet() CSet {
	m := CSet{
		m: make(map[string]map[string]string),
	}
	return m
}

// Add reason "fv" action "kickvote", reaction "voteban"
func (cs *CSet) Add(reason, action, reaction string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if cs.m[reason] == nil {
		cs.m[reason] = make(map[string]string)
	}
	cs.m[reason][action] = reaction
}

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

	return csvutil.Marshal(list)
}
