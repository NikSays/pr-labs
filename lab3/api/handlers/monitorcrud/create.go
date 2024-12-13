package monitorcrud

import (
	"encoding/json"
	"log"
	"net/http"

	"communicator/connections/db"
)

type createReqBody struct {
	Name     string  `json:"name"`
	PriceMDL float64 `json:"price_mdl"`
	PriceEUR float64 `json:"price_eur"`
	Warranty int32   `json:"warranty"`
}

func (g *HandlerGroup) create(rw http.ResponseWriter, req *http.Request) {
	var body createReqBody
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(rw, "Invalid json", http.StatusBadRequest)
		log.Print(err)
		return
	}

	if body.Warranty < 0 {
		http.Error(rw, "Warranty must be > 0", http.StatusBadRequest)
		return
	}

	err = g.Database.AddMonitor(req.Context(), db.AddMonitorParams{
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
