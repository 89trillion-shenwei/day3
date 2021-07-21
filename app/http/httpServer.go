package http

import "day3/internal/router"

func Start() {
	r := router.GiftRouter()

	r.Run("localhost:8000")
}
