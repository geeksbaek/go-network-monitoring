package main

import (
	"net"
	"sort"
	"time"

	"github.com/pivotal-golang/bytefmt"
)

var (
	localhost = []byte{172, 31}
)

// for Flow structure
type Endpoints struct {
	Src string
	Dst string
}

type Traffic struct {
	Total uint64
	// Inbound  uint64
	// Outbound uint64

	FormattedTotal string
}

type Flow struct {
	*Endpoints
	*Traffic
}

type Flows []*Flow

func (fs *Flows) Append(e *Endpoints, t *Traffic) {
	*fs = append(*fs, &Flow{e, t})
}

// for Statistic structure
type Stat struct {
	TimeFrom time.Time
	TimeTo   time.Time
	Talkers
	*Traffic
}

type Talker struct {
	Name string
	*Traffic
}

type Talkers []*Talker

func (fs Flows) ToStat() *Stat {

	total := uint64(0)

	tsMap := map[string]*Traffic{}
	for _, f := range fs {
		total += f.Total
		get(tsMap, me(f.Endpoints)).Total += f.Total
		get(tsMap, me(f.Endpoints)).FormattedTotal = bytefmt.ByteSize(get(tsMap, me(f.Endpoints)).Total)
	}

	// fmt.Printf("Flows Length is %d / Talkers Length is %d.\n", len(fs), len(tsMap))

	ts := Talkers{}
	for k, v := range tsMap {
		ts = append(ts, &Talker{k, v})
	}

	sort.Sort(sort.Reverse(ts))

	return &Stat{
		TimeFrom: TimeFrom,
		TimeTo:   time.Now(),
		Talkers:  ts,
		Traffic:  &Traffic{total, bytefmt.ByteSize(total)},
	}
}

func get(m map[string]*Traffic, k string) *Traffic {
	if m[k] == nil {
		m[k] = &Traffic{}
	}
	return m[k]
}

func me(e *Endpoints) string {
	dst := net.ParseIP(e.Dst)
	for i := range localhost {
		if localhost[i] != dst[i] {
			return e.Src
		}
	}
	return e.Dst
}

// For Stats sorting
func (t Talkers) Len() int           { return len(t) }
func (t Talkers) Less(i, j int) bool { return t[i].Total < t[j].Total }
func (t Talkers) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
