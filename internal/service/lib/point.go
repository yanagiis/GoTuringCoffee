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
	X    *float64  `json:"x"`
	Y    *float64  `json:"y"`
	Z    *float64  `json:"z"`
	E    *float64  `json:"e"`
	E1   *float64  `json:"e1"`
	E2   *float64  `json:"e2"`
	F    *float64  `json:"f"`
	T    *float64  `json:"t"`
	Time *float64  `json:"time"`
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
	return math.Pow(sum, 0.5)
}
