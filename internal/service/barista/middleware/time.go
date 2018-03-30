package middleware

import "github.com/yanagiis/GoTuringCoffee/internal/service/lib"

type TimeMiddleware struct {
	pos lib.Point
***REMOVED***

func NewTimeMiddleware(***REMOVED*** *TimeMiddleware {
	x := float64(0***REMOVED***
	y := float64(0***REMOVED***
	z := float64(0***REMOVED***
	return &TimeMiddleware{
		pos: lib.Point{
			X: &x,
			Y: &y,
			Z: &z,
	***REMOVED***,
***REMOVED***
***REMOVED***

func (m *TimeMiddleware***REMOVED*** SetPos(p *lib.Point***REMOVED*** {
	if p.X != nil {
		x := *p.X
		m.pos.X = &x
***REMOVED***
	if p.Y != nil {
		y := *p.Y
		m.pos.Y = &y
***REMOVED***
	if p.Z != nil {
		z := *p.Z
		m.pos.Z = &z
***REMOVED***
***REMOVED***

func (m *TimeMiddleware***REMOVED*** Transform(p *lib.Point***REMOVED*** {
	if p.Time == nil {
		return
***REMOVED***
	if p.X == nil && p.Y == nil && p.Z == nil {
		p.Time = p.F
***REMOVED*** else {
		distance := m.pos.CalcDistance(p***REMOVED***
		*p.Time = distance * 60 / *p.F
***REMOVED***
	m.SetPos(p***REMOVED***
***REMOVED***
