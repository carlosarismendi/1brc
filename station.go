package main

import "sync"

type Station struct {
	Name  string
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

type StationMap struct {
	mutex sync.Mutex
	m     map[string]*Station
}

func NewStationMap() *StationMap {
	return &StationMap{
		m: make(map[string]*Station),
	}
}

func (sm *StationMap) Add(name string, temperature float64) {
	// sm.mutex.Lock()
	// defer sm.mutex.Unlock()

	s, ok := sm.m[name]
	if ok {
		if temperature < s.Min {
			s.Min = temperature
		}

		if temperature > s.Max {
			s.Max = temperature
		}

		s.Sum += temperature
		s.Count++
	} else {
		sm.m[name] = &Station{
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
