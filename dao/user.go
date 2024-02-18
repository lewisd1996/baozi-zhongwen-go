package dao

import (
	. "github.com/go-jet/jet/v2/postgres"
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

func (dao *Dao) GetUserById(id string) (User, error) {
	var user User
	stmt := table.User.SELECT(table.User.AllColumns).WHERE(table.User.ID.EQ(UUID(uuid.MustParse(id)))).LIMIT(1)
	err := stmt.Query(dao.DB, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}
