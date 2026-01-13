package models

import "time"

var subscrTimeLayout = "01-2006"

type Subscription struct {
	Id                int        `json:"id" db:"id"`
	ServiceName       string     `json:"service_name" db:"service_name"`
	Price             int        `json:"price" db:"price"`
	UserId            string     `json:"user_id" db:"user_id"`
	StartDate         time.Time  `json:"-" db:"start_date"`
	EndDate           *time.Time `json:"-" db:"end_date"`
	StartDateFormated string     `json:"start_date" db:"-"`
	EndDateFormated   *string    `json:"end_date" db:"-"`
}

func (s *Subscription) Format() {
	s.StartDateFormated = s.StartDate.Format(subscrTimeLayout)

	if s.EndDate != nil {
		endStr := s.EndDate.Format(subscrTimeLayout)
		s.EndDateFormated = &endStr
	}
}

func (s *Subscription) Parse() error {
	var err error
	s.StartDate, err = time.Parse(subscrTimeLayout, s.StartDateFormated)
	if err != nil {
		return err
	}

	if s.EndDateFormated != nil {
		end, err := time.Parse(subscrTimeLayout, *s.EndDateFormated)
		if err != nil {
			return err
		}
		s.EndDate = &end
	}

	return nil
}
