package moviecrud

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"communicator/connections/db"
)

type updateReqBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Rating      uint8  `json:"rating"`
}

func (g *HandlerGroup) update(rw http.ResponseWriter, req *http.Request) {
	idString := req.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(rw, "ID in path must be an int", http.StatusBadRequest)
		return
	}

	var body updateReqBody
	err = json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(rw, "Invalid json", http.StatusBadRequest)
		log.Print(err)
		return
	}

	if body.Rating < 0 || body.Rating > 5 {
		http.Error(rw, "Rating must be between 0 and 5", http.StatusBadRequest)
		return
	}

	err = g.Database.UpdateMovie(req.Context(), db.UpdateMovieParams{
		ID:          int32(id),
		Name:        body.Name,
		Description: body.Description,
		Rating:      int32(body.Rating),
	})
	if err != nil {
		http.Error(rw, "Failed saving to the database", http.StatusInternalServerError)
		log.Print(err)
		return
	}
}
