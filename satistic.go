package main

import (
	"sync/atomic"
)

//SatisticType is the type of satistic value
type SatisticType int

const (
	//ConnectionCount Connection number counts
	ConnectionCount SatisticType = iota + 1
	//StreamCount Stream number counts
	StreamCount
)

//Satistic internal Satistic type and its value
type Satistic struct {
	SatisticType SatisticType
	value        int64
}

//NodeStatisticMessage Statistic message
type NodeStatisticMessage struct {
	Dataset []Satistic
}

// Load atomically loads the wrapped value.
func (s *Satistic) Load() int64 {
	return atomic.LoadInt64(&s.value)
}

// Add atomically adds to the wrapped int64 and returns the new value.
func (s *Satistic) Add(n int64) int64 {
	return atomic.AddInt64(&s.value, n)
}

// Sub atomically subtracts from the wrapped int64 and returns the new value.
func (s *Satistic) Sub(n int64) int64 {
	return atomic.AddInt64(&s.value, -n)
}

//GetStatistic return specific type of satistic
func (m *NodeStatisticMessage) GetStatistic(stype SatisticType) (Satistic, error) {
	for _, sat := range m.Dataset {
		if sat.SatisticType == stype {
			return sat, nil
		}
	}
	return Satistic{}, errNotFound
}

//Get return all statistics
func (m *NodeStatisticMessage) Get() []Satistic {
	return m.Dataset
}
