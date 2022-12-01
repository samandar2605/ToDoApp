package postgres_test

import (
	"testing"

	"github.com/ToDoApp/storage/repo"
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

func createUser(t *testing.T) *repo.User {
	user, err := strg.User().Create(&repo.User{
		FirstName:   faker.Name(),
		LastName:    faker.LastName(),
		PhoneNumber: faker.Phonenumber(),
		Email:       faker.Email(),
		ImageUrl:    faker.URL(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)
	return user
}

func deleteUser(id int, t *testing.T) {
	err := strg.User().Delete(id)
	require.NoError(t, err)
}

func TestGetUser(t *testing.T) {
	u := createUser(t)
	usr, err := strg.User().GetById(u.Id)
	require.NoError(t, err)
	require.NotEmpty(t, usr)

	deleteUser(usr.Id, t)
}

func TestCreateUser(t *testing.T) {
	createUser(t)
}

func TestUpdateUser(t *testing.T) {
	u := createUser(t)

	u.FirstName = faker.Name()
	u.LastName = faker.LastName()
	u.PhoneNumber = faker.Phonenumber()
	u.Email = faker.Email()
	u.ImageUrl = faker.URL()

	user, err := strg.User().Update(u)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	deleteUser(user.Id, t)
}

func TestDeleteUser(t *testing.T) {
	u := createUser(t)
	deleteUser(u.Id, t)
}

func TestGetAllUser(t *testing.T) {
	u := createUser(t)
	n, _ := faker.RandomInt(100)
	_, err := strg.User().GetAll(&repo.GetAllUsersParams{
		Page:  n[0],
		Limit: n[1],
	})

	require.NoError(t, err)
	deleteUser(u.Id, t)
}
