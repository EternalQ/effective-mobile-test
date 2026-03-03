package models_test

import (
	"testing"
	"time"

	"github.com/EternalQ/effective-mobile-test/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestSubscription_Parse(t *testing.T) {
	tests := []struct {
		name  string
		start string
		end   string

		wantStart bool
		wantEnd   bool

		wantErr bool
	}{
		{
			name:      "normal",
			start:     "03-2026",
			end:       "05-2026",
			wantStart: true,
			wantEnd:   true,
			wantErr:   false,
		},
		{
			name:      "no start",
			start:     "",
			end:       "05-2026",
			wantStart: false,
			wantEnd:   true,
			wantErr:   false,
		},
		{
			name:      "no end",
			start:     "03-2026",
			end:       "",
			wantStart: true,
			wantEnd:   false,
			wantErr:   false,
		},
		{
			name:      "zero end",
			start:     "03-2026",
			end:       "0",
			wantStart: true,
			wantEnd:   false,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &models.Subscription{
				StartDateFormated: tt.start,
				EndDateFormated:   tt.end,
			}
			gotErr := s.Parse()
			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.Nil(t, gotErr)
			}

			if tt.wantStart {
				want, err := time.Parse(models.SubscrTimeLayout, s.StartDateFormated)
				assert.Nil(t, err)
				assert.Equal(t, want, s.StartDate)
			}

			if tt.wantEnd {
				want, err := time.Parse(models.SubscrTimeLayout, s.EndDateFormated)
				assert.Nil(t, err)
				assert.Equal(t, want, *s.EndDate)
			}
		})
	}
}
