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
	m     map[string]Station
}

func NewStationMap() *StationMap {
	return &StationMap{
		m: make(map[string]Station),
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
		sm.m[s.Name] = s
	} else {
		sm.m[station.Name] = station
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
		sm.m[s.Name] = s
	} else {
		sm.m[name] = Station{
			Name:  name,
			Min:   temperature,
			Max:   temperature,
			Sum:   temperature,
			Count: 1,
		}
	}
}
