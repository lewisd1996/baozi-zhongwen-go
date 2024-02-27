//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Deck = newDeckTable("public", "deck", "")

type deckTable struct {
	postgres.Table

	// Columns
	ID          postgres.ColumnString
	CreatedAt   postgres.ColumnTimestampz
	UpdatedAt   postgres.ColumnTimestampz
	Name        postgres.ColumnString
	Description postgres.ColumnString
	OwnerID     postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type DeckTable struct {
	deckTable

	EXCLUDED deckTable
}

// AS creates new DeckTable with assigned alias
func (a DeckTable) AS(alias string) *DeckTable {
	return newDeckTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new DeckTable with assigned schema name
func (a DeckTable) FromSchema(schemaName string) *DeckTable {
	return newDeckTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new DeckTable with assigned table prefix
func (a DeckTable) WithPrefix(prefix string) *DeckTable {
	return newDeckTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new DeckTable with assigned table suffix
func (a DeckTable) WithSuffix(suffix string) *DeckTable {
	return newDeckTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newDeckTable(schemaName, tableName, alias string) *DeckTable {
	return &DeckTable{
		deckTable: newDeckTableImpl(schemaName, tableName, alias),
		EXCLUDED:  newDeckTableImpl("", "excluded", ""),
	}
}

func newDeckTableImpl(schemaName, tableName, alias string) deckTable {
	var (
		IDColumn          = postgres.StringColumn("id")
		CreatedAtColumn   = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn   = postgres.TimestampzColumn("updated_at")
		NameColumn        = postgres.StringColumn("name")
		DescriptionColumn = postgres.StringColumn("description")
		OwnerIDColumn     = postgres.StringColumn("owner_id")
		allColumns        = postgres.ColumnList{IDColumn, CreatedAtColumn, UpdatedAtColumn, NameColumn, DescriptionColumn, OwnerIDColumn}
		mutableColumns    = postgres.ColumnList{CreatedAtColumn, UpdatedAtColumn, NameColumn, DescriptionColumn, OwnerIDColumn}
	)

	return deckTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:          IDColumn,
		CreatedAt:   CreatedAtColumn,
		UpdatedAt:   UpdatedAtColumn,
		Name:        NameColumn,
		Description: DescriptionColumn,
		OwnerID:     OwnerIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}