package inmem

import (
	"github.com/hashicorp/go-memdb"
)

func New() (*memdb.MemDB, error) {
	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"person": {
				Name: "person",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "id"},
					},
					"username": {
						Name:    "age",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "username"},
					},
				},
			},
		},
	}

	return memdb.NewMemDB(schema)
}
