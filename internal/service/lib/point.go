package lib

***REMOVED***
	"math"
***REMOVED***

type Point struct {
	X    *float64 `json:"x,omitempty"`
	Y    *float64 `json:"y,omitempty"`
	Z    *float64 `json:"z,omitempty"`
	E    *float64 `json:"e,omitempty"`
	E1   *float64 `json:"e1,omitempty"`
	E2   *float64 `json:"e2,omitempty"`
	F    *float64 `json:"f,omitempty"`
	T    *float64 `json:"t,omitempty"`
	Time *float64 `json:"time,omitempty"`
***REMOVED***

func (p *Point***REMOVED*** CalcDistance(other *Point***REMOVED*** float64 {
	var sum float64 = 0
	if p.X != nil && other.X != nil {
		sum += math.Pow(*p.X-*other.X, 2***REMOVED***
***REMOVED***
	if p.Y != nil && other.Y != nil {
		sum += math.Pow(*p.Y-*other.Y, 2***REMOVED***
***REMOVED***
	if p.Z != nil && other.Z != nil {
		sum += math.Pow(*p.Z-*other.Z, 2***REMOVED***
***REMOVED***
	return sum
***REMOVED***
