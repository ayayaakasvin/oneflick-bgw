package models

import "time"

type EventStats struct {
	EventUUID         string    `json:"event_uuid"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	StartingTime      time.Time `json:"starting_time"`
	EndingTime        time.Time `json:"ending_time"`
	Status            string    `json:"status"`
	Capacity          int       `json:"capacity"`
	ImageURL          string    `json:"image_url"`
	CategoryName      string    `json:"category_name"`
	OrganizerUsername string    `json:"organizer_username"`
	TotalTicketsSold  int       `json:"total_tickets_sold"`
	Rank              int       `json:"rank"`

	// Derived field — computed in Go, not stored in DB
	FillRate float64 `json:"fill_rate"`
}
