package lib

import (
	"math"
)

type PointType int

const (
	PointT PointType = iota
	WaitT
	MixT
	HomeT
)

type Point struct {
	Type PointType `json:"type"`
	X    *float64  `json:"x,omitempty"`
	Y    *float64  `json:"y,omitempty"`
	Z    *float64  `json:"z,omitempty"`
	E    *float64  `json:"e,omitempty"`
	E1   *float64  `json:"e1,omitempty"`
	E2   *float64  `json:"e2,omitempty"`
	F    *float64  `json:"f,omitempty"`
	T    *float64  `json:"t,omitempty"`
	Time *float64  `json:"time,omitempty"`
}

func (p *Point) CalcDistance(other *Point) float64 {
	var sum float64 = 0
	if p.X != nil && other.X != nil {
		sum += math.Pow(*p.X-*other.X, 2)
	}
	if p.Y != nil && other.Y != nil {
		sum += math.Pow(*p.Y-*other.Y, 2)
	}
	if p.Z != nil && other.Z != nil {
		sum += math.Pow(*p.Z-*other.Z, 2)
	}
	return sum
}
