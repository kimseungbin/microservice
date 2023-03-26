package http

import (
	"encoding/json"
	"errors"
	"log"
	"microservices/rating/internal/controller/rating"
	model "microservices/rating/pkg"
	"net/http"
	"strconv"
)

// Handler defines a rating service controller.
type Handler struct {
	ctrl *rating.Controller
}

// New create a new rating service HTTP handler.
func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// Handle handles PUT and GET /rating requests.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	id := model.RecordID(r.FormValue("id"))
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	recordType := model.RecordType(r.FormValue("type"))
	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		v, err := h.ctrl.GetAggregatedRating(r.Context(), id, recordType)
		if err != nil {
			if errors.Is(err, rating.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
		if err := json.NewEncoder(w).Encode(v); err != nil {
			log.Printf("Response encode error: %v\n", err)
			return
		}
	case http.MethodPost:
		userID := model.UserID(r.FormValue("userId"))
		v, err := strconv.ParseFloat(r.FormValue("value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := h.ctrl.PutRating(
			r.Context(),
			id,
			recordType,
			&model.Rating{UserID: userID, Value: model.RatingValue(v)},
		); err != nil {
			log.Printf("Repository put error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
