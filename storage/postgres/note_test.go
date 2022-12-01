package postgres_test

import (
	"fmt"
	"testing"

	"github.com/ToDoApp/storage/repo"
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

func createNote(t *testing.T) *repo.Note {
	n, _ := faker.RandomInt(10)
	Note, err := strg.Note().Create(&repo.Note{
		UserId:      n[1],
		Title:       faker.Sentence(),
		Description: faker.Sentence(),
	})
	fmt.Println(Note.UserId)
	require.NoError(t, err)
	require.NotEmpty(t, Note)
	return Note
}

func deleteNote(id int, t *testing.T) {
	err := strg.Note().Delete(id)
	require.NoError(t, err)
}

func TestGetNote(t *testing.T) {
	n := createNote(t)
	note, err := strg.Note().GetById(n.Id)
	require.NoError(t, err)
	require.NotEmpty(t, note)

	deleteNote(note.Id, t)
}

func TestCreateNote(t *testing.T) {
	createNote(t)
}

func TestUpdateNote(t *testing.T) {
	n := createNote(t)

	//UserId is connected with user. So, I gave value
	n.UserId = 1

	n.Title = faker.Sentence()
	n.Description = faker.Sentence()

	Note, err := strg.Note().Update(n)
	require.NoError(t, err)
	require.NotEmpty(t, Note)

	deleteNote(Note.Id, t)
}

func TestDeleteNote(t *testing.T) {
	u := createNote(t)
	deleteNote(u.Id, t)
}

func TestGetAllNote(t *testing.T) {
	u := createNote(t)
	n, _ := faker.RandomInt(100)
	_, err := strg.Note().GetAll(&repo.GetAllNotesParams{
		Page:  n[0],
		Limit: n[0],
	})

	require.NoError(t, err)
	deleteNote(u.Id, t)
}
