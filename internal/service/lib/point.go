package lib

***REMOVED***
	"bytes"
	"errors"
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

func (p *Point***REMOVED*** ToGCode(***REMOVED*** (string, error***REMOVED*** {
	if p.X == nil && p.Y == nil && p.Z == nil {
		return "", errors.New("no x, y, and z"***REMOVED***
***REMOVED***

	var buffer bytes.Buffer
	buffer.WriteString("G1"***REMOVED***
	if p.X != nil {
		buffer.WriteString(fmt.Sprintf(" X%0.5f", *p.X***REMOVED******REMOVED***
***REMOVED***
	if p.Y != nil {
		buffer.WriteString(fmt.Sprintf(" Y%0.5f", *p.Y***REMOVED******REMOVED***
***REMOVED***
	if p.Z != nil {
		buffer.WriteString(fmt.Sprintf(" Z%0.5f", *p.Z***REMOVED******REMOVED***
***REMOVED***
	buffer.WriteString(fmt.Sprintf(" F%0.5f", *p.F***REMOVED******REMOVED***
	return buffer.String(***REMOVED***, nil
***REMOVED***

func (p *Point***REMOVED*** ToHCode(***REMOVED*** (string, error***REMOVED*** {
	if p.Time == nil {
		return "", errors.New("no time"***REMOVED***
***REMOVED***

	var buffer bytes.Buffer
	buffer.WriteString("H"***REMOVED***
	if p.E1 != nil && *p.E1 != 0 {
		buffer.WriteString(fmt.Sprintf(" E0 %05f", *p.E1***REMOVED******REMOVED***
***REMOVED***
	if p.E2 != nil && *p.E2 != 0 {
		buffer.WriteString(fmt.Sprintf(" E1 %05f", *p.E2***REMOVED******REMOVED***
***REMOVED***
	buffer.WriteString(fmt.Sprintf(" T%0.5f", *p.Time***REMOVED******REMOVED***
	return buffer.String(***REMOVED***, nil
***REMOVED***
