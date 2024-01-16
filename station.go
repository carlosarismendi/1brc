package main

import (
	"fmt"
	"hash/maphash"
	"sort"
	"strings"
)

type Station struct {
	Name  string
	Min   float64
	Max   float64
	Sum   float64
	Count int64
}

func (s Station) String() string {
	mean := s.Sum / float64(s.Count)
	return fmt.Sprintf("%s=%.1f/%.1f/%.1f", s.Name, round(s.Min), mean, round(s.Max))
}

var seed = maphash.MakeSeed()

type StationMap struct {
	m map[uint64]*Station
}

func NewStationMap() *StationMap {
	return &StationMap{
		m: make(map[uint64]*Station),
	}
}

func (sm *StationMap) Add(name string, nameBytes []byte, temperature float64) {
	n := maphash.Bytes(seed, nameBytes)

	if s, ok := sm.m[n]; ok {
		if temperature < s.Min {
			s.Min = temperature
		}

		if temperature > s.Max {
			s.Max = temperature
		}

		s.Sum += temperature
		s.Count++
	} else {
		sm.m[n] = &Station{
			Name:  name,
			Min:   temperature,
			Max:   temperature,
			Sum:   temperature,
			Count: 1,
		}
	}
}

func (sm *StationMap) Merge(smap *StationMap) {
	left := sm.m
	right := smap.m
	for k := range right {
		r := right[k]
		if l, ok := left[k]; ok {
			l.Min = min(l.Min, r.Min)
			l.Max = max(l.Max, r.Max)
			l.Sum += r.Sum
			l.Count += r.Count
			continue
		}

		left[k] = r
		delete(right, k)
	}
}

func (sm *StationMap) String() string {
	stations := make([]*Station, 0, len(sm.m))
	for k := range sm.m {
		stations = append(stations, sm.m[k])
		delete(sm.m, k)
	}

	sort.Slice(stations, func(i, j int) bool {
		return stations[i].Name < stations[j].Name
	})

	var sb strings.Builder
	sep := ""
	for _, v := range stations {
		sb.WriteString(fmt.Sprintf("%s%s", sep, v.String()))
		sep = ", "
	}
	return sb.String()
}
