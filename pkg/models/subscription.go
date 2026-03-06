package models

import "time"

var SubscrTimeLayout = "01-2006"

type Subscription struct {
	Id                 int        `json:"id" db:"id"`
	ServiceName        string     `json:"service_name" db:"service_name"`
	Price              int        `json:"price" db:"price"`
	UserId             string     `json:"user_id" db:"user_id"`
	StartDate          time.Time  `json:"-" db:"start_date"`
	EndDate            *time.Time `json:"-" db:"end_date"`
	StartDateFormatted string     `json:"start_date" db:"-"`
	EndDateFormatted   string     `json:"end_date" db:"-"` //omitempty?
}

func (s *Subscription) Format() {
	s.StartDateFormatted = s.StartDate.Format(SubscrTimeLayout)

	if s.EndDate != nil {
		s.EndDateFormatted = s.EndDate.Format(SubscrTimeLayout)
	}
}

func (s *Subscription) Parse() error {
	var err error
	if s.StartDateFormatted != "" {
		s.StartDate, err = time.Parse(SubscrTimeLayout, s.StartDateFormatted)
		if err != nil {
			return err
		}
	}

	if s.EndDateFormatted != "" && s.EndDateFormatted != "0" {
		end, err := time.Parse(SubscrTimeLayout, s.EndDateFormatted)
		if err != nil {
			return err
		}
		s.EndDate = &end
	}

	return nil
}
