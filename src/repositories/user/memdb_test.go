package user

import (
	"github.com/stretchr/testify/assert"
	"golang-server/src/infrastructure/inmem"
	"testing"
)

func TestUserRepositories(t *testing.T) {

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
			err := repo.Create(user.username, "")
			assert.NoErrorf(t, err, "create user failed %s %d", user.username, iR)
			userInDb, err := repo.Read(nil, &user.username)
			assert.NoErrorf(t, err, "read user by username failed %s %d", user.username, iR)
			assert.Equal(t, userInDb.Username, user.username)
			userInDb, err = repo.Read(&userInDb.Id, nil)
			err = repo.Delete(userInDb.Id)
			assert.NoErrorf(t, err, "delete user by id failed %s %d userInDb %#v", user.username, iR, userInDb)

			userInDb, err = repo.Read(&userInDb.Id, nil)
			assert.ErrorIsf(t, err, NotFoundError, "read after delete user by id failed %s %d userInDb %#v", user.username, iR, userInDb)
		}
	}
}
