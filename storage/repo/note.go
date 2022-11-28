package repo

import "time"

type Note struct {
	Id          int
	UserId      int
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type NoteStorageI interface {
	Create(note *Note) (*Note, error)
	GetById(id int) (*Note, error)
	GetAll(params *GetAllNotesParams) (*GetAllNotesResult, error)
	Update(note *Note) (*Note, error)
	Delete(id int) error
}

type GetAllNotesParams struct {
	Limit      int
	Page       int
	SortByDate string
	Search     string
}

type GetAllNotesResult struct {
	Notes []*Note
	Count int
}
