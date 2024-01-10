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
	mutex sync.RWMutex
	m     map[string]*Station
}

func NewStationMap() *StationMap {
	return &StationMap{
		m: make(map[string]*Station),
	}
}

func (sm *StationMap) Update(station Station) {
	s, ok := sm.m[station.Name]
	if ok {
		temperature := station.Min
		if temperature < s.Min {
			s.Min = temperature
		}

		if temperature > s.Max {
			s.Max = temperature
		}

		s.Sum += temperature
		s.Count++
	} else {
		sm.m[station.Name] = &station
	}
}
