//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package view

import (
	"github.com/go-jet/jet/v2/postgres"
)

var VCardWithLastReviewed = newVCardWithLastReviewedTable("public", "v_card_with_last_reviewed", "")

type vCardWithLastReviewedTable struct {
	postgres.Table

	// Columns
	ID             postgres.ColumnString
	CreatedAt      postgres.ColumnTimestampz
	UpdatedAt      postgres.ColumnTimestampz
	DeckID         postgres.ColumnString
	Content        postgres.ColumnString
	Translation    postgres.ColumnString
	LastReviewedAt postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type VCardWithLastReviewedTable struct {
	vCardWithLastReviewedTable

	EXCLUDED vCardWithLastReviewedTable
}

// AS creates new VCardWithLastReviewedTable with assigned alias
func (a VCardWithLastReviewedTable) AS(alias string) *VCardWithLastReviewedTable {
	return newVCardWithLastReviewedTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new VCardWithLastReviewedTable with assigned schema name
func (a VCardWithLastReviewedTable) FromSchema(schemaName string) *VCardWithLastReviewedTable {
	return newVCardWithLastReviewedTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new VCardWithLastReviewedTable with assigned table prefix
func (a VCardWithLastReviewedTable) WithPrefix(prefix string) *VCardWithLastReviewedTable {
	return newVCardWithLastReviewedTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new VCardWithLastReviewedTable with assigned table suffix
func (a VCardWithLastReviewedTable) WithSuffix(suffix string) *VCardWithLastReviewedTable {
	return newVCardWithLastReviewedTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newVCardWithLastReviewedTable(schemaName, tableName, alias string) *VCardWithLastReviewedTable {
	return &VCardWithLastReviewedTable{
		vCardWithLastReviewedTable: newVCardWithLastReviewedTableImpl(schemaName, tableName, alias),
		EXCLUDED:                   newVCardWithLastReviewedTableImpl("", "excluded", ""),
	}
}

func newVCardWithLastReviewedTableImpl(schemaName, tableName, alias string) vCardWithLastReviewedTable {
	var (
		IDColumn             = postgres.StringColumn("id")
		CreatedAtColumn      = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn      = postgres.TimestampzColumn("updated_at")
		DeckIDColumn         = postgres.StringColumn("deck_id")
		ContentColumn        = postgres.StringColumn("content")
		TranslationColumn    = postgres.StringColumn("translation")
		LastReviewedAtColumn = postgres.TimestampzColumn("last_reviewed_at")
		allColumns           = postgres.ColumnList{IDColumn, CreatedAtColumn, UpdatedAtColumn, DeckIDColumn, ContentColumn, TranslationColumn, LastReviewedAtColumn}
		mutableColumns       = postgres.ColumnList{IDColumn, CreatedAtColumn, UpdatedAtColumn, DeckIDColumn, ContentColumn, TranslationColumn, LastReviewedAtColumn}
	)

	return vCardWithLastReviewedTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:             IDColumn,
		CreatedAt:      CreatedAtColumn,
		UpdatedAt:      UpdatedAtColumn,
		DeckID:         DeckIDColumn,
		Content:        ContentColumn,
		Translation:    TranslationColumn,
		LastReviewedAt: LastReviewedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
