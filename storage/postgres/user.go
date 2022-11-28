package postgres

import (
	"fmt"
	"time"

	"github.com/ToDoApp/storage/repo"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	DB *sqlx.DB
}

func NewUser(db *sqlx.DB) repo.UserStorageI {
	return &userRepo{
		DB: db,
	}
}

func (ur *userRepo) Create(user *repo.User) (*repo.User, error) {
	query := `
		insert into users(
			first_name,
			last_name,
			phone_number,
			email,
			password,
			image_url)
		values($1,$2,$3,$4,$5,$6)
		returning id,created_at
	`
	row := ur.DB.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.PhoneNumber,
		user.Email,
		user.Password,
		user.ImageUrl,
	)

	if err := row.Scan(&user.Id, &user.CreatedAt); err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepo) GetById(id int) (*repo.User, error) {
	var (
		user repo.User
	)

	query := `
		select 
			id,
			first_name,
			last_name,
			phone_number,
			email,
			password,
			image_url,
			created_at
		from users
		where id=$1 and deleted_at is null
	`
	row := ur.DB.QueryRow(query, id)

	if err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Email,
		&user.Password,
		&user.ImageUrl,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepo) GetAll(params *repo.GetAllUsersParams) (*repo.GetAllUsersResult, error) {
	result := repo.GetAllUsersResult{
		Users: make([]*repo.User, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)

	filter := ""
	if params.Search != "" {
		filter += " WHERE first_name ilike '%" + params.Search + "%' or last_name ilike '%" + params.Search + "%' or email ilike '%" + params.Search + "%' "
	}

	query := `
		SELECT
			id,
			first_name,
			last_name,
			phone_number,
			email,
			password,
			image_url,
			created_at
		FROM users
		` + filter + `
		ORDER BY created_at ` + params.SortByDate + limit

	rows, err := ur.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user repo.User

		err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.PhoneNumber,
			&user.Email,
			&user.Password,
			&user.ImageUrl,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Users = append(result.Users, &user)
	}
	fmt.Println(result)
	queryCount := `SELECT count(1) FROM users ` + filter
	err = ur.DB.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil

}

func (ur *userRepo) Update(user *repo.User) (*repo.User, error) {
	query := `
		update users set 
			first_name=$1,
			last_name=$2,
			phone_number=$3,
			email=$4,
			image_url=$5,
			updated_at=$6
		where id=$7 and deleted_at IS NULL
		returning updated_at
	`
	row := ur.DB.QueryRow(query,
		user.FirstName,
		user.LastName,
		user.PhoneNumber,
		user.Email,
		user.ImageUrl,
		time.Now(),
		user.Id,
	)

	if err := row.Scan(
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepo) Delete(id int) error {
	query := `
		update users set 
			deleted_at=$1
		where id=$2
	`
	_, err := ur.DB.Exec(
		query,
		time.Now(),
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepo) GetByEmail(email string) (*repo.User, error) {
	var result repo.User

	query := `
		SELECT
			id,
			first_name,
			last_name,
			phone_number,
			email,
			password,
			image_url,
			created_at
		FROM users
		WHERE email=$1
	`

	row := ur.DB.QueryRow(query, email)
	err := row.Scan(
		&result.Id,
		&result.FirstName,
		&result.LastName,
		&result.PhoneNumber,
		&result.Email,
		&result.Password,
		&result.ImageUrl,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) UpdatePassword(req *repo.UpdatePassword) error {
	query := `UPDATE users SET password=$1 WHERE id=$2`

	_, err := ur.DB.Exec(query, req.Password, req.UserID)
	if err != nil {
		return err
	}

	return nil
}
