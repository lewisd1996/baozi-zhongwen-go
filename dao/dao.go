package dao

import "database/sql"

type Dao struct {
	DB *sql.DB
}

/* -------------------------------------------------------------------------- */
/*                                    Init                                    */
/* -------------------------------------------------------------------------- */

func NewDao(db *sql.DB) *Dao {
	return &Dao{
		DB: db,
	}
}
