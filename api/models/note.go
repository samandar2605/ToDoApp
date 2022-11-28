package models

import "time"

type Note struct {
	Id          int       `json:"id"`
	UserId      int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type CreateNote struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GetAllNotesResponse struct {
	Notes []*Note `json:"users"`
	Count int  `json:"count"`
}

type GetAllNotesParams struct {
	Limit      int    `json:"limit"`
	Page       int    `json:"page"`
	Search     string `json:"search"`
	SortByDate string `json:"sort_by_date"`
}
