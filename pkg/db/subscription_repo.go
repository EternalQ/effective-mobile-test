package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/EternalQ/effective-mobile-test/pkg/models"
	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("entity not found")

type SubscriptionRepo struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewSubscriptionRepo(db *sqlx.DB, log *slog.Logger) *SubscriptionRepo {
	return &SubscriptionRepo{
		db,
		log.With(slog.String("where", "db/SubscriptionRepo")),
	}
}

var createSubscription = `
INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;`

func (r *SubscriptionRepo) Create(s *models.Subscription) error {
	err := r.db.Get(s, createSubscription, s.ServiceName, s.Price, s.UserId, s.StartDate, s.EndDate)
	if err != nil {
		r.log.Error("Error while creating entity",
			slog.String("err", err.Error()),
			slog.String("method", "Create"),
		)
		return err
	}

	return nil
}

var readSubscription = `
SELECT * 
FROM subscriptions 
WHERE id = $1`

func (r *SubscriptionRepo) Read(id int) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.Get(&subscription, readSubscription, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		r.log.Error("Error while getting entity",
			slog.String("err", err.Error()),
			slog.String("method", "Read"),
		)
		return nil, err
	}
	return &subscription, nil
}

// var updateSubsription = `
// UPDATE subscriptions
// SET service_name = :service_name, price = :price, user_id = :user_id, start_date = :start_date, end_date = :end_date
// WHERE id = :id`

func (r *SubscriptionRepo) Update(subscription *models.Subscription) error {
	fields := []string{}
	if subscription.UserId != "" {
		fields = append(fields, "user_id = :user_id")
	}
	if subscription.ServiceName != "" {
		fields = append(fields, "service_name = :service_name")
	}
	if subscription.Price != 0 {
		fields = append(fields, "price = :price")
	}
	if !subscription.StartDate.IsZero() {
		fields = append(fields, "start_date = :start_date")
	}
	if subscription.EndDateFormated == "0" {
		fields = append(fields, "end_date = NULL")
	} else if subscription.EndDate != nil && !subscription.EndDate.IsZero() {
		fields = append(fields, "end_date = :end_date")
	}

	// if len(fields) == 0 {
	// 	return err
	// }
	query := fmt.Sprintf("UPDATE subscriptions SET %s WHERE id = :id", strings.Join(fields, ", "))
	r.log.Debug("Update query", slog.String("string", query))

	res, err := r.db.NamedExec(query, subscription)
	if err != nil {
		r.log.Error("Error while updating entity",
			slog.String("err", err.Error()),
			slog.String("method", "Update"),
		)
		return err
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		r.log.Debug("Nothing updated",
			slog.String("method", "Update"),
		)
		return ErrNotFound
	}

	return nil
}

var deleteSubscription = `
DELETE FROM subscriptions 
WHERE id = $1`

func (r *SubscriptionRepo) Delete(id int) error {
	res, err := r.db.Exec(deleteSubscription, id)
	if err != nil {
		r.log.Error("Error while creating entity",
			slog.String("err", err.Error()),
			slog.String("method", "Delete"),
		)
		return err
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		r.log.Debug("Nothing deleted",
			slog.String("method", "Delete"),
		)
		return ErrNotFound
	}

	return nil
}

// var listSubscription = `
// SELECT *
// FROM subscriptions
// WHERE 1=1`

func (r *SubscriptionRepo) List(filter *models.Subscription) ([]*models.Subscription, error) {
	var subscriptions []*models.Subscription
	query := "SELECT * FROM subscriptions"

	if filter != nil {
		conditions := []string{}

		if filter.UserId != "" {
			conditions = append(conditions, "user_id = :user_id")
		}
		if filter.ServiceName != "" {
			conditions = append(conditions, "service_name = :service_name")
		}
		if !filter.StartDate.IsZero() {
			conditions = append(conditions, "start_date >= :start_date")
		}
		if filter.EndDate != nil && !filter.EndDate.IsZero() {
			conditions = append(conditions, "end_date <= :end_date")
		}

		if len(conditions) > 0 {
			query += " WHERE " + strings.Join(conditions, " AND ")
		}
	} else {
		filter = &models.Subscription{}
	}
	r.log.Debug("Select query", slog.String("string", query))

	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, query, filter)
	r.log.Debug("Prepared query", slog.String("string", query), slog.Any("args", args))
	if err != nil {
		r.log.Error("Error while preparing quary",
			slog.String("err", err.Error()),
			slog.String("method", "List"),
		)
		return nil, err
	}
	err = r.db.Select(&subscriptions, query, args...)
	if err != nil {
		r.log.Error("Error while listing entity",
			slog.String("err", err.Error()),
			slog.String("method", "List"),
		)
		return nil, err
	}

	return subscriptions, nil
}
