package outlet

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jackc/pgx/v5"
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

func (r *Repo) NewOutlet(txa pgx.Tx, name string, address string, postal string, officialLinks []string) (id int64, err error) {
	if txa == nil {
		return 0, fmt.Errorf("tx is nil")
	}

	defer func() {
		if err != nil {
			txa.Rollback(context.Background())
		}
	}()

	row := txa.QueryRow(context.Background(), "insert into outlet(name,address,postal_code, official_links) values ($1,$2,$3,$4) returning id;", name, address, postal, officialLinks)
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo) NewMenuItem(txa pgx.Tx, name string) (id int64, err error) {
	if name == "" {
		return 0, fmt.Errorf("menu item name is empty")
	}
	if txa == nil {
		return 0, fmt.Errorf("tx is nil")
	}

	defer func() {
		if err != nil {
			txa.Rollback(context.Background())
		}
	}()

	row := txa.QueryRow(context.Background(), "insert into menu_item(name) values ($1) returning id;", name)

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo) NewStallMenuItem(txa pgx.Tx, menuItemId int64, outletId int64) (id int64, err error) {
	if txa == nil {
		txa, err = r.conn.Begin(context.Background())
		if err != nil {
			return 0, err
		}
	}

	defer func() {
		if err != nil {
			txa.Rollback(context.Background())
		} else {
			txa.Commit(context.Background())
		}
	}()

	row := txa.QueryRow(context.Background(), "insert into outlet_menu(outlet_id, menu_item_id) values ($1,$2) returning id;", outletId, menuItemId)

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo) NewOutletWithMenu(outletName string, address string, postal string, officialLinks []string, menuItems []string) (err error) {
	tx, err := r.conn.Begin(context.Background())
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	outletId, err := r.NewOutlet(tx, outletName, address, postal, officialLinks)
	if err != nil {
		return err
	}

	if len(menuItems) > 1 {
		return fmt.Errorf("outlet %s has more than one menu item. server not implemented", outletName)
	}
	if len(menuItems) == 0 {
		return fmt.Errorf("outlet %s has no item server", outletName)
	}

	itemId, err := r.NewMenuItem(tx, menuItems[0])
	if err != nil {
		return err
	}

	_, err = r.NewStallMenuItem(tx, itemId, outletId)
	if err != nil {
		return err
	}

	tx.Commit(context.Background())
	return nil
}

func New(conn *pgxpool.Pool) *Repo {
	return &Repo{conn: conn}
}
