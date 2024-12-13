package monitorcrud

import (
	"encoding/json"
	"log"
	"net/http"

	"communicator/connections/db"
)

type readReqBody struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"page_size"`
}

func (g *HandlerGroup) read(rw http.ResponseWriter, req *http.Request) {
	var body readReqBody
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(rw, "Invalid json", http.StatusBadRequest)
		log.Print(err)
		return
	}

	movies, err := g.Database.AllMonitors(req.Context(), db.AllMonitorsParams{
		Limit:  body.PageSize,
		Offset: body.Page * body.PageSize,
	})
	if err != nil {
		http.Error(rw, "Failed reading from the database", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	moviesJSON, err := json.Marshal(movies)
	if err != nil {
		http.Error(rw, "Error building JSON", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	_, err = rw.Write(moviesJSON)
	if err != nil {
		http.Error(rw, "Error writing response", http.StatusInternalServerError)
		log.Print(err)
		return
	}
}
