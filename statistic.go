package main

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/pivotal-golang/bytefmt"
)

type pair struct {
	Key   string
	Value uint64
}

type pairList []pair

func (p pairList) Len() int           { return len(p) }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func printSortedStatisticMap(stat map[string]uint64) {
	p := make(pairList, len(stat))

	i := 0
	for k, v := range stat {
		p[i] = pair{k, v}
		i++
	}

	sort.Sort(sort.Reverse(p))
	var b bytes.Buffer

	for i, k := range p {
		if i >= 20 {
			break
		}
		fmt.Fprintf(&b, "%15s = %s\n", k.Key, bytefmt.ByteSize(k.Value))
	}
	fmt.Print(b.String())
}
