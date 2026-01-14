package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/EternalQ/effective-mobile-test/pkg/db"
	"github.com/EternalQ/effective-mobile-test/pkg/models"
	"github.com/EternalQ/effective-mobile-test/pkg/service"
	"github.com/gorilla/mux"
)

type Server struct {
	log      *slog.Logger
	subsServ *service.SubscriptionService
}

func StartServer(log *slog.Logger, subsServ *service.SubscriptionService, r *mux.Router) {
	s := &Server{
		log.With(slog.String("where", "api/Server")),
		subsServ,
	}

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/subscriptions", s.createSubsription).Methods("POST")
	api.HandleFunc("/subscriptions", s.listSubsription).Methods("get")
	api.HandleFunc("/subscriptions/{id}", s.readSubsription).Methods("Get")
	api.HandleFunc("/subscriptions/{id}", s.updateSubscription).Methods("PATCH")
	api.HandleFunc("/subscriptions/{id}", s.deleteSubscription).Methods("DELETE")
	api.HandleFunc("/subscriptions/calc", s.calculateSubscription).Methods("POST")
}

func (s *Server) handleError(w http.ResponseWriter, r *http.Request, err string, code int) {
	s.log.Error("Error while handling request",
		slog.String("err", err),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)
	http.Error(w, err, code)
}

// @Summary Create a new subscription
// @Description Creates a new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.Subscription true "Subscription details"
// @Success 201 {int} int "ID of the created subscription"
// @Failure 400 {string} string "Invalid input"
// @Router /subscriptions [post]
func (s *Server) createSubsription(w http.ResponseWriter, r *http.Request) {
	s.log.Info("Handling POST request to /api/subscriptions")

	var sub *models.Subscription
	err := json.NewDecoder(r.Body).Decode(&sub)
	if err != nil {
		s.handleError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	if err := sub.Parse(); err != nil {
		s.log.Error("Error while parsing subscription",
			slog.String("err", err.Error()),
			slog.String("source", "models/Subscription.Parse"),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.subsServ.Create(sub); err != nil {
		s.handleError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	s.log.Info("Subscription created", slog.Int("id", sub.Id))
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"id": sub.Id})
}

// @Summary List all subscriptions
// @Description Lists all subscriptions
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {array} models.Subscription "List of subscriptions"
// @Router /subscriptions [get]
func (s *Server) listSubsription(w http.ResponseWriter, r *http.Request) {
	s.log.Info("Handling GET request to /api/subscriptions")

	subs, err := s.subsServ.ListAll()
	if err != nil {
		s.handleError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, sub := range subs {
		sub.Format()
	}

	s.log.Info("Subscriptions listed", slog.Int("count", len(subs)))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subs)
}

// @Summary Read a subscription by ID
// @Description Reads a subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} models.Subscription "Subscription details"
// @Failure 404 {string} string "Subscription not found"
// @Router /subscriptions/{id} [get]
func (s *Server) readSubsription(w http.ResponseWriter, r *http.Request) {
	s.log.Info("Handling GET request to /api/subscriptions/{id}")

	vars := mux.Vars(r)
	s.log.Debug("GET /api/subscriptions/{id}", slog.String("id", vars["id"]))

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.handleError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	sub, err := s.subsServ.Read(id)
	if err == db.ErrNotFound {
		s.handleError(w, r, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		s.handleError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	sub.Format()

	s.log.Info("Subscription readed", slog.Int("id", sub.Id))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

// @Summary Update a subscription by ID
// @Description Updates a subscription by ID. Needs at least 1 field to update. Send `"end_date": "0"` to set null.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param subscription body models.Subscription true "Accepted fields of Subscription: `service_name`, `price`, `user_id`, `start_date`, `end_date`"
// @Success 200 {object} models.Subscription "Updated subscription details"
// @Failure 400 {string} string "Invalid input"
// @Router /subscriptions/{id} [patch]
func (s *Server) updateSubscription(w http.ResponseWriter, r *http.Request) {
	s.log.Info("Handling PATCH request to /api/subscriptions/{id}")

	vars := mux.Vars(r)
	s.log.Debug("PATCH /api/subscriptions/{id}", slog.String("id", vars["id"]))

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.handleError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	var sub *models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		s.handleError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	if sub.UserId == "" &&
		sub.ServiceName == "" &&
		sub.Price == 0 &&
		sub.StartDateFormated == "" &&
		sub.EndDateFormated == "" {
		s.handleError(w, r, "Empty fields", http.StatusBadRequest)
		return
	}

	if err := sub.Parse(); err != nil {
		s.log.Error("Error while parsing subscription",
			slog.String("err", err.Error()),
			slog.String("source", "models/Subscription.Parse"),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sub.Id = id
	err = s.subsServ.Update(sub)
	if err == db.ErrNotFound {
		s.handleError(w, r, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		s.handleError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	s.log.Info("Subscription updated", slog.Int("id", sub.Id))
	w.WriteHeader(http.StatusOK)
}

// @Summary Delete a subscription by ID
// @Description Deletes a subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 204 {string} string "Subscription deleted"
// @Failure 404 {string} string "Subscription not found"
// @Router /subscriptions/{id} [delete]
func (s *Server) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	s.log.Info("Handling DELETE request to /api/subscriptions/{id}")

	vars := mux.Vars(r)
	s.log.Debug("DELETE /api/subscriptions/{id}", slog.String("id", vars["id"]))

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.handleError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.subsServ.Delete(id)
	if err == db.ErrNotFound {
		s.handleError(w, r, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		s.handleError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	s.log.Info("Subscription deleted", slog.Int("id", id))
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Calculate subscription price
// @Description Calculates the subscription price based on a filter. Looks for all subscriptions with given `user_id` and `service_name` between `start_date` and `end_date`. Fields may be omitted but needs at least 1.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param filter body models.Subscription true "Filter for subscription calculation"
// @Success 200 {int} int "Calculated price"
// @Failure 400 {string} string "Invalid input"
// @Router /subscriptions/calc [post]
func (s *Server) calculateSubscription(w http.ResponseWriter, r *http.Request) {
	s.log.Info("Handling POST request to /api/subscriptions/calc")

	var filter *models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		s.handleError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	if err := filter.Parse(); err != nil {
		s.log.Error("Error while parsing subscription",
			slog.String("err", err.Error()),
			slog.String("source", "models/Subscription.Parse"),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	price, err := s.subsServ.CalculatePrice(filter)
	if err != nil {
		s.handleError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	s.log.Info("Price calculated", slog.Int("price", price))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"price": price})
}
