package models

import "time"

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Train struct {
	ID            int
	Name          string
	DepartureTime time.Time
	ArrivalTime   time.Time
	From          string
	To            string
	TotalSeats    int
	Fare          int
	CreatedAt     time.Time
}
