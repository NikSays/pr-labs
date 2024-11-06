package moviecrud

import (
	"encoding/json"
	"log"
	"net/http"

	"communicator/connections/db"
)

type createReqBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Rating      uint8  `json:"rating"`
}

func (g *HandlerGroup) create(rw http.ResponseWriter, req *http.Request) {
	var body createReqBody
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(rw, "Invalid json", http.StatusBadRequest)
		log.Print(err)
		return
	}

	if body.Rating < 0 || body.Rating > 5 {
		http.Error(rw, "Rating must be between 0 and 5", http.StatusBadRequest)
		return
	}

	err = g.Database.AddMovie(req.Context(), db.AddMovieParams{
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
