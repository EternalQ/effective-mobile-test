package service_test

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/EternalQ/effective-mobile-test/pkg/models"
	"github.com/EternalQ/effective-mobile-test/pkg/service"
	"github.com/stretchr/testify/assert"
)

var ErrNotImplemented = errors.New("not implemented")

type MockRepo struct {
	listFn func(*models.Subscription) ([]*models.Subscription, error)
}

func (m *MockRepo) Create(*models.Subscription) error      { return ErrNotImplemented }
func (m *MockRepo) Read(int) (*models.Subscription, error) { return nil, ErrNotImplemented }
func (m *MockRepo) Update(*models.Subscription) error      { return ErrNotImplemented }
func (m *MockRepo) Delete(int) error                       { return ErrNotImplemented }

func (m *MockRepo) List(f *models.Subscription) ([]*models.Subscription, error) {
	return m.listFn(f)
}

func TestSubscriptionService_CalculatePrice(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	t.Run("Price calculation", func(t *testing.T) {
		m := &MockRepo{
			func(s *models.Subscription) ([]*models.Subscription, error) {
				return []*models.Subscription{
					{Price: 100},
					{Price: 200},
					{Price: 300},
				}, nil
			},
		}

		ss := service.NewSubscriptionService(m, logger)
		price, err := ss.CalculatePrice(nil)

		assert.Nil(t, err)
		assert.Equal(t, 600, price)
	})

	t.Run("db error", func(t *testing.T) {
		m := &MockRepo{
			func(s *models.Subscription) ([]*models.Subscription, error) {
				return nil, errors.New("db error")
			},
		}

		ss := service.NewSubscriptionService(m, logger)
		price, err := ss.CalculatePrice(nil)

		assert.Error(t, err)
		assert.Equal(t, -1, price)
	})
}
