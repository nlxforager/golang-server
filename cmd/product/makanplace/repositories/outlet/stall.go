package outlet

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

func (r *Repo) updateOutlet(tx pgx.Tx, outletId *int64, name string, address string, postal string, officialLinks []string) (err error) {
	if tx == nil {
		return fmt.Errorf("tx is nil")
	}
	if outletId == nil {
		return fmt.Errorf("[updateOutlet] outletId is nil")
	}

	_, err = tx.Exec(context.Background(), "UPDATE outlet set name=$1, address=$2, postal_code=$3, official_links=$4 where id=$5", name, address, postal, officialLinks, *outletId)
	return err
}
func (r *Repo) newOutlet(tx pgx.Tx, name string, address string, postal string, officialLinks []string) (id int64, err error) {
	if tx == nil {
		return 0, fmt.Errorf("tx is nil")
	}

	row := tx.QueryRow(context.Background(), "insert into outlet(name,address,postal_code, official_links) values ($1,$2,$3,$4) returning id;", name, address, postal, officialLinks)
	err = row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("[newOutlet] %w", err)
	}
	return id, nil
}

func (r *Repo) newMenuItems(txa pgx.Tx, names []string) (_ids []int64, err error) {
	if txa == nil {
		return []int64{}, fmt.Errorf("tx is nil")
	}

	args := []interface{}{}
	placeholders := []string{}

	for i, name := range names {
		if name == "" {
			return []int64{}, fmt.Errorf("some menu item name is empty")
		}
		placeholders = append(placeholders, fmt.Sprintf("($%d)", i+1))
		args = append(args, name)
	}

	query := fmt.Sprintf("INSERT INTO menu_item(name) VALUES %s RETURNING id;", strings.Join(placeholders, ", "))

	rows, err := txa.Query(context.Background(), query, args...)
	if err != nil {
		return []int64{}, fmt.Errorf("[NewMenuItem] %w", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var _id int64
		err = rows.Scan(&_id)
		if err != nil {
			return []int64{}, fmt.Errorf("[NewMenuItem] %w", err)
		}
		ids = append(ids, _id)
	}
	if err != nil {
		return []int64{}, fmt.Errorf("[NewMenuItem] %w", err)
	}
	return ids, nil
}

func (r *Repo) NewStallMenuItem(txa pgx.Tx, menuItemIds []int64, outletId int64, replace bool) (id int64, err error) {
	if txa == nil {
		return 0, fmt.Errorf("tx is nil")
	}

	if replace {
		_, err = txa.Exec(context.Background(), "DELETE FROM outlet_menu where outlet_id = $1 ", outletId)
		if err != nil {
			return 0, fmt.Errorf("[NewStallMenuItem] %w", err)
		}
	}

	args := []interface{}{}
	placeholders := []string{}

	for i, name := range menuItemIds {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d)", i*2+1, i*2+2))
		args = append(args, outletId, name)
	}

	query := fmt.Sprintf("insert into outlet_menu(outlet_id, menu_item_id) values %s returning id;", strings.Join(placeholders, ", "))
	log.Printf("query %s", query)
	rows, err := txa.Query(context.Background(), query, args...)
	defer rows.Close()
	if err != nil {
		return 0, fmt.Errorf("[NewStallMenuItem] %w", err)
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo) UpdateOutletWithMenu(outletId *int64, outletName string, address string, postal string, officialLinks []string, reviewLinks []string, menuItems []string) (err error) {
	tx, err := r.conn.Begin(context.Background())
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	err = r.updateOutlet(tx, outletId, outletName, address, postal, officialLinks)
	if err != nil {
		return err
	}
	{
		if len(menuItems) == 0 {
			return fmt.Errorf("edit outlet request %s: menu required", outletName)
		}

		itemIds, err := r.newMenuItems(tx, menuItems)
		if err != nil {
			return err
		}

		_, err = r.NewStallMenuItem(tx, itemIds, *outletId, true)
		if err != nil {
			return err
		}
	}
	err = r.replaceReviewLinks(tx, outletId, reviewLinks)
	if err != nil {
		return err
	}
	return tx.Commit(context.Background())
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

	outletId, err := r.newOutlet(tx, outletName, address, postal, officialLinks)
	if err != nil {
		return err
	}

	{
		if len(menuItems) == 0 {
			return fmt.Errorf("new outlet request %s: menu required", outletName)
		}

		itemIds, err := r.newMenuItems(tx, menuItems)
		if err != nil {
			return err
		}

		_, err = r.NewStallMenuItem(tx, itemIds, outletId, false)
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

	MenuItems json.RawMessage `json:"menu_items"` // raw JSON
}

func (r *Repo) GetOutlets(postalCode *string, id *int) ([]Outlet, error) {
	var selectCriteria strings.Builder
	if postalCode != nil {
		selectCriteria.WriteString(fmt.Sprintf(" where postal_code='%s'", *postalCode))
	}

	if id != nil {
		selectCriteria.WriteString(fmt.Sprintf(" where outlet.id=%d", *id))
	}

	rows, err := r.conn.Query(context.Background(), "select outlet.id, outlet.name, address, postal_code, official_links, latlong, COALESCE(array_agg(owr.link) filter (where owr.link is not null), '{}'), (json_agg(json_build_object('id', mi.id, 'name', mi.name)))  from outlet left join public.outlet_web_reviews owr on outlet.id = owr.outlet_id LEFT JOIN outlet_menu om on om.outlet_id = outlet.id LEFT JOIN menu_item mi on mi.id = om.menu_item_id"+
		selectCriteria.String()+
		" group by outlet.id;")
	if err != nil {
		return nil, fmt.Errorf("GetOutlets: %w", err)
	}
	var outlets []Outlet
	for rows.Next() {

		var o Outlet

		rows.Scan(&o.Id, &o.Name, &o.Address, &o.PostalCode, &o.OfficialLinks, &o.LatLong, &o.ReviewLinks, &o.MenuItems)
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
	return err
}

func (r *Repo) replaceReviewLinks(tx pgx.Tx, outletId *int64, links []string) error {
	if tx == nil {
		return fmt.Errorf("tx is nil")
	}
	if outletId == nil {
		return fmt.Errorf("outletId is nil")
	}
	_, err := tx.Exec(context.Background(), "DELETE FROM outlet_web_reviews where outlet_id=$1", *outletId)
	if err != nil {
		return err
	}

	_, err = tx.CopyFrom(context.Background(), pgx.Identifier{"outlet_web_reviews"}, []string{"outlet_id", "link"}, pgx.CopyFromSlice(len(links), func(i int) ([]any, error) {
		return []any{outletId, links[i]}, nil
	}))
	return err
}

func New(conn *pgxpool.Pool) *Repo {
	return &Repo{conn: conn}
}
