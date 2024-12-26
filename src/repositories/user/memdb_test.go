package user

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"golang-server/src/infrastructure/inmem"
)

func TestUserCRUD(t *testing.T) {
	type User struct {
		username string
	}

	usersToCreate := []User{{"abc"}}

	inMemDb, err := inmem.New()
	if err != nil {
		panic(err)
	}

	repoInMem, err := NewMemDbUserRepository(inMemDb)
	if err != nil {
		panic(err)
	}

	repos := []UserRepository{repoInMem}

	for iR, repo := range repos {
		for _, user := range usersToCreate {
			password := "12345"
			wrongPassword := "asvas"
			err := repo.Create(user.username, password)
			assert.NoErrorf(t, err, "create user failed %s %d", user.username, iR)
			userInDb, err := repo.Read(nil, &user.username, nil)
			assert.NoErrorf(t, err, "read user by username failed %s %d", user.username, iR)

			{
				_, err := repo.Read(nil, &user.username, &password)
				assert.NoErrorf(t, err, "read user by username and password failed %s %d", user.username, iR)
			}
			_, err = repo.Read(nil, &user.username, &wrongPassword)
			assert.Error(t, err, "read user by username should fail %s %d", user.username, iR)
			assert.Equal(t, userInDb.Username, user.username)
			userInDb, err = repo.Read(&userInDb.Id, nil, &password)
			err = repo.Delete(userInDb.Id)
			assert.NoErrorf(t, err, "delete user by id failed %s %d userInDb %#v", user.username, iR, userInDb)

			userInDb, err = repo.Read(&userInDb.Id, nil, &password)
			assert.ErrorIsf(t, err, NotFoundError, "read after delete user by id failed %s %d userInDb %#v", user.username, iR, userInDb)
		}
	}
}
