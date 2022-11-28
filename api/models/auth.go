package models

import "time"

type RegisterRequest struct {
	FirstName   string `json:"first_name" binding:"required,min=2,max=30"`
	LastName    string `json:"last_name" binding:"required,min=2,max=30"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6,max=16"`
	PhoneNumber string `json:"phone_number"`
	ImageUrl    string `json:"image_url"`
}

type AuthResponse struct {
	Id          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	ImageUrl    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"update_at"`
	DeletedAt   time.Time `json:"deleted_at"`
	AccessToken string    `json:"access_token"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=16"`
}

type VerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}
