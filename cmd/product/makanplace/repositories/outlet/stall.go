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
		return 0, fmt.Errorf("tx is nil")
	}

	row := txa.QueryRow(context.Background(), "insert into outlet_menu(outlet_id, menu_item_id) values ($1,$2) returning id;", outletId, menuItemId)

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo) NewOutletWithMenu(outletName string, address string, postal string, officialLinks []string, reviewLinks []string, menuItems []string) (err error) {
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

	{
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
	}

	err = r.addReviewLinks(tx, outletId, reviewLinks)
	if err != nil {
		return err
	}

	tx.Commit(context.Background())
	return nil
}

type Outlet struct {
	LatLong       []string
	Name          string
	Address       string
	PostalCode    string
	OfficialLinks []string
	ReviewLinks   []string
	Id            int64
}

func (r *Repo) GetOutlets() ([]Outlet, error) {
	rows, err := r.conn.Query(context.Background(), "select outlet.id, name, address, postal_code, official_links, latlong, array_agg(owr.link) filter (where owr.link is not null) from outlet left join public.outlet_web_reviews owr on outlet.id = owr.outlet_id group by outlet.id;")
	if err != nil {
		return nil, err
	}
	var outlets []Outlet
	for rows.Next() {

		var o Outlet

		rows.Scan(&o.Id, &o.Name, &o.Address, &o.PostalCode, &o.OfficialLinks, &o.LatLong, &o.ReviewLinks)
		outlets = append(outlets, o)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return outlets, nil
}

func (r *Repo) SetLatLong(postal string, latitude string, longitude string) error {
	_, err := r.conn.Exec(context.Background(), "update outlet set latlong=$1 where postal_code=$2", []string{latitude, longitude}, postal)
	return err
}

func (r *Repo) addReviewLinks(tx pgx.Tx, outletId int64, links []string) error {
	if tx == nil {
		return fmt.Errorf("tx is nil")
	}

	_, err := tx.CopyFrom(context.Background(), pgx.Identifier{"outlet_web_reviews"}, []string{"outlet_id", "link"}, pgx.CopyFromSlice(len(links), func(i int) ([]any, error) {
		return []any{outletId, links[i]}, nil
	}))
	if err != nil {
		return err
	}
	return nil
}

func New(conn *pgxpool.Pool) *Repo {
	return &Repo{conn: conn}
}
