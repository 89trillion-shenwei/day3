package http

import "day3/internal/router"

func Start() {
	r := router.SetStrRouter()

	r.Run("localhost:8000")
}
