package repo

import "time"

type User struct {
	Id          int
	FirstName   string
	LastName    string
	Password    string
	PhoneNumber string
	Email       string
	ImageUrl    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type UserStorageI interface {
	Create(user *User) (*User, error)
	GetById(id int) (*User, error)
	GetAll(params *GetAllUsersParams) (*GetAllUsersResult, error)
	Update(user *User) (*User, error)
	Delete(id int) error
	UpdatePassword(req *UpdatePassword) error
	GetByEmail(email string) (*User, error)
}

type GetAllUsersParams struct {
	Limit      int
	Page       int
	Search     string
	SortByDate string
}

type GetAllUsersResult struct {
	Users []*User
	Count int32
}

type UpdatePassword struct {
	UserID   int
	Password string
}
