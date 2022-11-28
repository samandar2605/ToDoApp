package storage

import (
	"github.com/ToDoApp/storage/postgres"
	"github.com/ToDoApp/storage/repo"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	User() repo.UserStorageI
	Note() repo.NoteStorageI
}

type storagePg struct {
	userRepo repo.UserStorageI
	noteRepo repo.NoteStorageI
}

func NewStoragePg(db *sqlx.DB) StorageI {
	return &storagePg{
		userRepo: postgres.NewUser(db),
		noteRepo: postgres.NewNote(db),
	}
}

func (s *storagePg) User() repo.UserStorageI {
	return s.userRepo
}

func (s *storagePg) Note() repo.NoteStorageI {
	return s.noteRepo
}
