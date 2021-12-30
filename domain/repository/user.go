package repository

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mine/just-projecting/domain/model"
	"github.com/mine/just-projecting/pkg/database"
)

type UserRepository interface {
	Find() ([]*model.User, error)
	FindOne() (*model.User, error)
	CountUserRepository(userID int64) (result string, err error)
}

type userRepository struct {
	// db connection, config etc
	db   *sqlx.DB
	dbMs *sqlx.DB
	// inject here
}

func NewUserRepository() *userRepository {
	return &userRepository{
		db: database.DB,
	}
}

func (r *userRepository) Find() ([]*model.User, error) {
	u := make([]*model.User, 0)
	u = append(u, &model.User{ID: 1, Name: "I'm First User", Age: 21})
	u = append(u, &model.User{ID: 2, Name: "I'm Second User", Age: 22})
	return u, nil
}

func (r *userRepository) FindOne() (*model.User, error) {
	return &model.User{
		ID:   1,
		Name: "I'm User",
		Age:  20,
	}, nil
}

func (r *userRepository) CountUserRepository(userID int64) (result string, err error) {
	err = r.db.QueryRow(`
			select count(*) 
			from users 
			where id = ?
	`, userID).Scan(&result)
	if err != nil {
		return result, errors.New("[User][mysql][CountUserRepository][userID=" + fmt.Sprint(userID) + "] failed to query countAllUsers")
	}

	return result, nil
}
