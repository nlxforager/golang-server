package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	conn *pgxpool.Pool
}

type BaseUser struct {
	Id        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}
type UserWithGmail struct {
	Id     int64    `db:"id"`
	Gmails []string `db:"gmails"`
}

func (r Repo) GetOrCreateUserByGmail(emails []string) (*UserWithGmail, error) {
	return r.getOrCreateUserByGmail(emails, 0)
}

// GetOrCreateUserByGmail guarantees a valid user on success
func (r Repo) getOrCreateUserByGmail(emails []string, recursed int) (*UserWithGmail, error) {
	if recursed > 1 {
		return nil, errors.New("getOrCreateUserByGmail: too many recursion")
	}
	if len(emails) == 0 {
		return nil, errors.New("no emails provided")
	}
	email := emails[0]

	// check if gmails is associated.
	log.Println("[GetOrCreateUserByGmail] check if gmails is associated...")
	row := r.conn.QueryRow(context.Background(), "with gmails1 as (select user_id, gmail from gmails where gmail = $1) select u.id id, array_agg(g.gmail) gmails from users u inner join gmails1 g on u.id = g.user_id group by u.id", email)

	var user UserWithGmail
	err := row.Scan(&user.Id, &user.Gmails)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Println("[GetOrCreateUserByGmail] gmails is not associated...")
		newUser, err := r.CreateUser()
		if err != nil {
			return nil, err
		}
		log.Printf("New user: %v\n", newUser)
		if err := r.AssociateGmail(newUser.Id, email); err != nil {
			return nil, err
		}

		return r.getOrCreateUserByGmail(emails, recursed+1)
	} else if err != nil {
		log.Printf("[GetOrCreateUserByGmail] general error %#v", err)
		return nil, err
	}
	return &user, nil
}

func (r Repo) CreateUser() (*BaseUser, error) {
	var user BaseUser
	row := r.conn.QueryRow(context.Background(), "insert into users DEFAULT VALUES RETURNING id, created_at;")
	err := row.Scan(&user.Id, &user.CreatedAt)

	if err != nil {
		log.Printf("[CreateUser] scan error %#v\n", err)
		return nil, err
	}
	return &user, nil
}

func (r Repo) AssociateGmail(userId int64, gmail string) error {
	_, err := r.conn.Exec(context.Background(), "insert into gmails(user_id, gmail) values (@user_id, @gmail) RETURNING *;", pgx.NamedArgs{
		"user_id": userId,
		"gmail":   gmail,
	})
	return err
}

func New(conn *pgxpool.Pool) *Repo {
	return &Repo{conn: conn}
}
