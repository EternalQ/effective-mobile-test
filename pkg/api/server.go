package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

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
	if err != nil {
		return
	}

	sub.Format()

	s.log.Info("Subscription readed", slog.Int("id", sub.Id))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

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

	if err := sub.Parse(); err != nil {
		s.log.Error("Error while parsing subscription",
			slog.String("err", err.Error()),
			slog.String("source", "models/Subscription.Parse"),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sub.Id = id
	if err := s.subsServ.Update(sub); err != nil {
		s.handleError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	s.log.Info("Subscription updated", slog.Int("id", sub.Id))
	w.WriteHeader(http.StatusOK)
}

func (s *Server) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	s.log.Info("Handling DELETE request to /api/subscriptions/{id}")

	vars := mux.Vars(r)
	s.log.Debug("DELETE /api/subscriptions/{id}", slog.String("id", vars["id"]))

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.handleError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.subsServ.Delete(id); err != nil {
		s.handleError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	s.log.Info("Subscription deleted", slog.Int("id", id))
	w.WriteHeader(http.StatusOK)
}

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
