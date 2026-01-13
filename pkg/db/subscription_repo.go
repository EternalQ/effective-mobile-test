package db

import (
	"log/slog"

	"github.com/EternalQ/effective-mobile-test/pkg/models"
	"github.com/jmoiron/sqlx"
)

type SubscriptionRepo struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewSubscriptionRepo(db *sqlx.DB, log *slog.Logger) *SubscriptionRepo {
	return &SubscriptionRepo{db, log}
}

var createSubscription = `
INSERT INTO subsriptions (service_name, price, user_id, start_date, end_date)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;`

func (r *SubscriptionRepo) Create(s *models.Subscription) error {
	err := r.db.Get(s, createSubscription, s.ServiceName, s.Price, s.UserId, s.StartDate, s.EndDate)
	if err != nil {
		return err
	}

	return nil
}

var readSubscription = `
SELECT * 
FROM subsriptions 
WHERE id = $1`

func (r *SubscriptionRepo) Read(id int) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.Get(&subscription, readSubscription, id)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

var updateSubsription = `
UPDATE subsriptions 
SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 
WHERE id = $6`

func (r *SubscriptionRepo) Update(subscription *models.Subscription) error {
	_, err := r.db.Exec(updateSubsription,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserId,
		subscription.StartDate,
		subscription.EndDate, subscription.Id)

	if err != nil {
		return err
	}
	return nil
}

var deleteSubscription = `
DELETE FROM subsriptions 
WHERE id = $1`

func (r *SubscriptionRepo) Delete(id int) error {
	_, err := r.db.Exec(deleteSubscription, id)
	if err != nil {
		return err
	}
	return nil
}

var listSubscription = `
SELECT * 
FROM subsriptions`

func (r *SubscriptionRepo) List() ([]*models.Subscription, error) {
	var subscriptions []*models.Subscription
	err := r.db.Select(&subscriptions, listSubscription)
	if err != nil {
		return nil, err
	}
	return subscriptions, nil
}
