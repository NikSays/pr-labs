package monitorcrud

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"communicator/connections/db"
)

type updateReqBody struct {
	Name     string  `json:"name"`
	PriceMDL float64 `json:"price_mdl"`
	PriceEUR float64 `json:"price_eur"`
	Warranty int32   `json:"warranty"`
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

	if body.Warranty < 0 {
		http.Error(rw, "Warranty must be > 0", http.StatusBadRequest)
		return
	}

	err = g.Database.UpdateMonitor(req.Context(), db.UpdateMonitorParams{
		ID:       int32(id),
		Name:     body.Name,
		PriceMdl: body.PriceMDL,
		PriceEur: body.PriceEUR,
		Warranty: body.Warranty,
	})
	if err != nil {
		http.Error(rw, "Failed saving to the database", http.StatusInternalServerError)
		log.Print(err)
		return
	}
}
