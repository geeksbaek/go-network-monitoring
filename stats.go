package main

import (
	"sort"
	"sync"
)

type Endpoints struct {
	Src string
	Dst string
}

type Traffic uint64

type Flow struct {
	Total Traffic
	Vars  map[Endpoints]Traffic
	Mutex *sync.RWMutex
}

type Stat struct {
	Endpoints
	Traffic
}

type Stats []*Stat

func (f *Flow) Add(e Endpoints, traffic Traffic) {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()
	f.Vars[e] += traffic
	f.Total += traffic
}

func (f *Flow) GetStats() Stats {
	f.Mutex.RLock()
	defer f.Mutex.RUnlock()
	s := make(Stats, 0, len(f.Vars))
	for endpoints, traffic := range f.Vars {
		s = append(s, &Stat{endpoints, traffic})
	}
	sort.Sort(sort.Reverse(s))
	return s
}

// For Stats sorting
func (s Stats) Len() int           { return len(s) }
func (s Stats) Less(i, j int) bool { return s[i].Traffic < s[j].Traffic }
func (s Stats) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
