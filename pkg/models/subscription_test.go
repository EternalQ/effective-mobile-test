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
		{
			name:      "invalid start",
			start:     "wrong",
			end:       "03-2026",
			wantStart: true,
			wantEnd:   false,
			wantErr:   true,
		},
		{
			name:      "invalid end",
			start:     "03-2026",
			end:       "wrong",
			wantStart: true,
			wantEnd:   false,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &models.Subscription{
				StartDateFormatted: tt.start,
				EndDateFormatted:   tt.end,
			}
			gotErr := s.Parse()
			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.Nil(t, gotErr)
			}

			if tt.wantStart {
				want, err := time.Parse(models.SubscrTimeLayout, s.StartDateFormatted)
				if !tt.wantErr {
					assert.Nil(t, err)
				}
				assert.Equal(t, want, s.StartDate)
			}

			if tt.wantEnd {
				want, err := time.Parse(models.SubscrTimeLayout, s.EndDateFormatted)
				if !tt.wantErr {
					assert.Nil(t, err)
				}
				assert.Equal(t, want, *s.EndDate)
			}
		})
	}
}

func TestSubscription_Format(t *testing.T) {
	tests := []struct {
		name string

		start time.Time
		end   time.Time

		wantStart string
		wantEnd   string
	}{
		{
			name:      "only start",
			start:     time.Date(2026, 3, 1, 0, 0, 0, 0, time.FixedZone("GMT", 3)),
			wantStart: "03-2026",
			wantEnd:   "",
		},
		{
			name:      "full",
			start:     time.Date(2026, 3, 1, 0, 0, 0, 0, time.FixedZone("GMT", 3)),
			end:       time.Date(2026, 6, 1, 0, 0, 0, 0, time.FixedZone("GMT", 3)),
			wantStart: "03-2026",
			wantEnd:   "06-2026",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &models.Subscription{
				StartDate: tt.start,
			}
			if !tt.end.IsZero() {
				s.EndDate = &tt.end
			}
			s.Format()
			assert.Equal(t, tt.wantStart, s.StartDateFormatted)
			assert.Equal(t, tt.wantEnd, s.EndDateFormatted)
		})
	}
}

func BenchmarkSubscription_Format(b *testing.B) {
	start := time.Date(2026, 3, 1, 0, 0, 0, 0, time.FixedZone("GMT", 3))
	end := time.Date(2026, 6, 1, 0, 0, 0, 0, time.FixedZone("GMT", 3))
	s := &models.Subscription{
		StartDate: start,
		EndDate:   &end,
	}

	for b.Loop() {
		s.Format()
	}
}
