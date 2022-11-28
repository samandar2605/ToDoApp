package models

import "time"

type User struct {
	Id          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Password    string    `json:"password"`
	Email       string    `json:"email"`
	ImageUrl    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"update_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type CreateUser struct {
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	ImageUrl    string `json:"image_url"`
}

type GetAllUsersResponse struct {
	Users []*User `json:"users"`
	Count int32   `json:"count"`
}

type GetAllUsersParams struct {
	Limit      int    `json:"limit" binding:"required"`
	Page       int    `json:"page" binding:"required"`
	Search     string `json:"search"`
	SortByDate string `json:"sort_by_date"`
}
