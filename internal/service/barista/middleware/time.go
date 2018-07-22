package middleware

import "github.com/yanagiis/GoTuringCoffee/internal/service/lib"

type TimeMiddleware struct {
	pos lib.Point
}

func NewTimeMiddleware() *TimeMiddleware {
	x := float64(0)
	y := float64(0)
	z := float64(0)
	return &TimeMiddleware{
		pos: lib.Point{
			X: &x,
			Y: &y,
			Z: &z,
		},
	}
}

func (m *TimeMiddleware) setPos(p *lib.Point) {
	if p.X != nil {
		x := *p.X
		m.pos.X = &x
	}
	if p.Y != nil {
		y := *p.Y
		m.pos.Y = &y
	}
	if p.Z != nil {
		z := *p.Z
		m.pos.Z = &z
	}
}

func (m *TimeMiddleware) Transform(p *lib.Point) {
	if p.Time == nil {
		return
	}
	if p.X == nil && p.Y == nil && p.Z == nil {
		p.Time = p.F
	} else {
		distance := m.pos.CalcDistance(p)
		*p.Time = distance * 60 / *p.F
	}
	m.setPos(p)
}

func (m *TimeMiddleware) Free() {

}
