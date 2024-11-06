package moviecrud

import (
	"net/http"

	"communicator/connections/db"
)

type HandlerGroup struct {
	Database db.Queries
}

func (g *HandlerGroup) Mux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", g.create)
	mux.HandleFunc("GET /", g.read)
	mux.HandleFunc("PUT /{id}", g.update)
	mux.HandleFunc("DELETE /{id}", g.update)

	return mux
}
