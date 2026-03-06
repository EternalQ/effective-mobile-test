package service

import (
	"log/slog"

	"github.com/EternalQ/effective-mobile-test/pkg/models"
)

type Repository interface {
	Create(*models.Subscription) error
	Read(int) (*models.Subscription, error)
	Update(*models.Subscription) error
	Delete(int) error
	List(*models.Subscription) ([]*models.Subscription, error)
}

type SubscriptionService struct {
	log           *slog.Logger
	subscriptions Repository
}

func NewSubscriptionService(subRepo Repository, log *slog.Logger) *SubscriptionService {
	return &SubscriptionService{
		log.With(slog.String("where", "service/SubscriptionService")),
		subRepo,
	}
}

func (ss *SubscriptionService) Create(s *models.Subscription) error {
	return ss.subscriptions.Create(s)
}

func (ss *SubscriptionService) Read(id int) (*models.Subscription, error) {
	return ss.subscriptions.Read(id)
}

func (ss *SubscriptionService) Update(s *models.Subscription) error {
	return ss.subscriptions.Update(s)
}

func (ss *SubscriptionService) Delete(id int) error {
	return ss.subscriptions.Delete(id)
}

func (ss *SubscriptionService) ListAll() ([]*models.Subscription, error) {
	return ss.subscriptions.List(nil)
}

func (ss *SubscriptionService) CalculatePrice(filter *models.Subscription) (int, error) {
	subs, err := ss.subscriptions.List(filter)
	if err != nil {
		ss.log.Error("Error while calculating price",
			slog.String("source", "db/SubcriptionRepo.List"),
			slog.String("method", "CalculatePrice"),
		)
		return -1, err
	}

	total := 0
	for _, s := range subs {
		total += s.Price
	}

	return total, nil
}
