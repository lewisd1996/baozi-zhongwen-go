package daos

import (
	"github.com/google/uuid"
	. "github.com/lewisd1996/baozi-zhongwen/sql/.jet/bz/public/model"
	"github.com/lewisd1996/baozi-zhongwen/sql/.jet/bz/public/table"
)

func (dao *Dao) CreateUser(email string, id uuid.UUID) error {
	user := User{
		ID:    id,
		Email: email,
	}
	stmt := table.User.INSERT(table.User.ID, table.User.Email).MODEL(user).RETURNING(table.User.AllColumns)
	_, err := stmt.Exec(dao.DB)

	if err != nil {
		println("Error creating user:", err.Error())
		return err
	}

	return nil
}
