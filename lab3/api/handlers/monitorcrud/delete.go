package monitorcrud

import (
	"net/http"
	"strconv"
)

func (g *HandlerGroup) delete(rw http.ResponseWriter, req *http.Request) {
	idString := req.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(rw, "ID in path must be an int", http.StatusBadRequest)
		return
	}

	err = g.Database.DeleteMonitor(req.Context(), int32(id))
	if err != nil {
		http.Error(rw, "Failed deleting from database", http.StatusInternalServerError)
	}
}
