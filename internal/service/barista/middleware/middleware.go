package middleware

import "github.com/yanagiis/GoTuringCoffee/internal/service/lib"

type Middleware interface {
	Transform(p *lib.Point)
	Free()
}
