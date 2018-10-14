package middleware

import "GoTuringCoffee/internal/service/lib"

type Middleware interface {
	Transform(p *lib.Point)
	Free()
}
