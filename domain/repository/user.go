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
	ListMerchantOmzet(date string, username string) (res []model.MerchantOmzet, err error)
	GetListMerchant(username string) (res []model.MerchantOmzet, err error)
	ListOutletOmzet(date string, username string) (res []model.Outlet, err error)
	GetListOutlet(username string) (res []model.Outlet, err error)
	GetDataByUsername(username string) (user model.Auth, err error)
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

func (r *userRepository) ListMerchantOmzet(date string, username string) (res []model.MerchantOmzet, err error) {
	var data model.MerchantOmzet

	rows, err := r.db.Queryx(`
	select 
		t.merchant_id,
		m.merchant_name,
		IF(sum(t.bill_total)=null, 0, sum(t.bill_total)) as omzet,
		DATE_FORMAT(t.created_at , '%Y-%m-%d') as created_at 
	from 
		transactions t
	left join 
		merchants m on m.id = t.merchant_id 
	left join 
		users u on u.id = m.user_id
	where 
		DATE_FORMAT(t.created_at , '%Y-%m-%d') = ?
	and
		u.user_name = ?
	group by t.merchant_id, created_at
	order by t.merchant_id
	`, date, username)
	if err != nil {
		return res, errors.New("[ListMerchantOmzet][While Queryx] Error = " + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&data)
		if err != nil {
			return res, errors.New(err.Error())
		}

		res = append(res, data)
	}

	return
}

func (r *userRepository) GetListMerchant(username string) (res []model.MerchantOmzet, err error) {
	var data model.MerchantOmzet

	rows, err := r.db.Queryx(`
		select m.id as merchant_id, m.merchant_name from merchants m
		left join users u on u.id = m.user_id
		where u.user_name = ?
	`, username)
	if err != nil {
		return res, errors.New("[ListMerchantOmzet][While Queryx] Error = " + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&data)
		if err != nil {
			return res, errors.New(err.Error())
		}

		res = append(res, data)
	}

	return
}

func (r *userRepository) ListOutletOmzet(date string, username string) (res []model.Outlet, err error) {
	var data model.Outlet

	rows, err := r.db.Queryx(`
		select 
			t.merchant_id,
			o.id as outlet_id ,
			m.merchant_name,
			o.outlet_name,
			IF(sum(t.bill_total)=null, 0, sum(t.bill_total)) as omzet,
			DATE_FORMAT(t.created_at , '%Y-%m-%d') as created_at 
		from 
			transactions t
		left join 
			merchants m on m.id = t.merchant_id
		left join
			outlets o on o.id = t.outlet_id 
		left join
			users u on u.id = m.user_id
		where 
			DATE_FORMAT(t.created_at , '%Y-%m-%d') = ?
		and
			u.user_name = ?
		group by t.merchant_id, o.id, created_at
		order by t.merchant_id
	`, date, username)
	if err != nil {
		return res, errors.New("[ListOutletOmzet][While Queryx] Error = " + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&data)
		if err != nil {
			return res, errors.New(err.Error())
		}

		res = append(res, data)
	}

	return res, nil
}

func (r *userRepository) GetListOutlet(username string) (res []model.Outlet, err error) {
	var data model.Outlet

	rows, err := r.db.Queryx(`
		select o.id as outlet_id, o.merchant_id, m.merchant_name ,outlet_name 
		from outlets o
		left join merchants m on m.id = o.merchant_id
		left join users u on u.id = m.user_id
		where u.user_name = ?
		order by o.merchant_id 
	`, username)
	if err != nil {
		return res, errors.New("[GetListOutlet][While Queryx] Error = " + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&data)
		if err != nil {
			return res, errors.New(err.Error())
		}

		res = append(res, data)
	}

	return
}

func (r *userRepository) GetDataByUsername(username string) (user model.Auth, err error) {
	err = r.db.QueryRowx(`
		select u.id as user_id,u.name,u.password,u.user_name, u.created_at, u.updated_at 
		from users u
		where u.user_name = ?
		limit 1
	`, username).StructScan(&user)
	if err != nil {
		return
	}

	fmt.Println(user)

	return
}
