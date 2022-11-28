package postgres

import (
	"fmt"
	"time"

	"github.com/ToDoApp/storage/repo"
	"github.com/jmoiron/sqlx"
)

type noteRepo struct {
	DB *sqlx.DB
}

func NewNote(db *sqlx.DB) repo.NoteStorageI {
	return &noteRepo{
		DB: db,
	}
}

func (nr *noteRepo) Create(note *repo.Note) (*repo.Note, error) {
	query := `
		insert into notes(
			user_id,
			title,
			description
		)
		values($1,$2,$3)
		returning id,created_at
	`

	row := nr.DB.QueryRow(
		query,
		note.UserId,
		note.Title,
		note.Description,
	)

	if err := row.Scan(&note.Id, &note.CreatedAt); err != nil {
		return nil, err
	}
	return note, nil
}

func (nr *noteRepo) GetById(id int) (*repo.Note, error) {
	var (
		note repo.Note
	)

	query := `
		select 
			id,
			user_id,
			title,
			description,
			created_at
		from notes
		where id=$1 and deleted_at is null
	`
	
	row := nr.DB.QueryRow(query, id)

	if err := row.Scan(
		&note.Id,
		&note.UserId,
		&note.Title,
		&note.Description,
		&note.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &note, nil
}

func (nr *noteRepo) GetAll(params *repo.GetAllNotesParams) (*repo.GetAllNotesResult, error) {
	result := repo.GetAllNotesResult{
		Notes: make([]*repo.Note, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)

	filter := ""
	if params.Search != "" {
		filter += " WHERE title ilike '%" + params.Search + "%' "
	}

	query := `
		SELECT
			id,
			user_id,
			title,
			description,
			created_at
		FROM notes
		` + filter + `
		ORDER BY created_at ` + params.SortByDate + limit

	rows, err := nr.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var note repo.Note

		err := rows.Scan(
			&note.Id,
			&note.UserId,
			&note.Title,
			&note.Description,
			&note.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Notes = append(result.Notes, &note)
	}

	queryCount := `SELECT count(1) FROM notes ` + filter
	err = nr.DB.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil

}

func (nr *noteRepo) Update(note *repo.Note) (*repo.Note, error) {
	query := `
		update notes set 
			user_id=$1,
			title=$2,
			description=$3,
			updated_at=$4
		where id=$5 and deleted_at IS NULL
		returning updated_at
	`
	row := nr.DB.QueryRow(query,
		note.UserId,
		note.Title,
		note.Description,
		time.Now(),
		note.Id,
	)

	if err := row.Scan(
		&note.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return note, nil
}

func (nr *noteRepo) Delete(id int) error {
	query := `
		update notes set 
			deleted_at=$1
		where id=$2
	`
	_, err := nr.DB.Exec(
		query,
		time.Now(),
		id,
	)
	if err != nil {
		return err
	}

	return nil
}
