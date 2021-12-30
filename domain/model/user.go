package model

import "time"

type Auth struct {
	ID        string    `json:"user_id" db:"user_id"`
	Name      string    `json:"name" validate:"required" db:"name"`
	Password  string    `json:"password" validate:"required" db:"password"`
	Username  string    `json:"user_name" db:"user_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type User struct {
	ID   int
	Name string
	Age  int
}

type Login struct {
	Username string `json:"user_name" sql:"user_name"`
	Password string `json:"password" validate:"required" sql:"password"`
}

type MerchantOmzet struct {
	MerchantID   int64   `json:"merchant_id" db:"merchant_id"`
	MerchantName string  `json:"merchant_name" db:"merchant_name"`
	Omzet        float64 `json:"omzet" db:"omzet"`
	CreatedAt    string  `json:"created_at" db:"created_at"`
}

type Outlet struct {
	MerchantID   int64   `json:"merchant_id" db:"merchant_id"`
	OutletID     int64   `json:"outlet_id" db:"outlet_id"`
	MerchantName string  `json:"merchant_name" db:"merchant_name"`
	OutletName   string  `json:"outlet_name" db:"outlet_name"`
	Omzet        float64 `json:"omzet" db:"omzet"`
	CreatedAt    string  `json:"created_at" db:"created_at"`
}

type LoginResponse struct {
	Token string
	User  Auth
}

type Token struct {
	Token string `json:"Authorization" db:"Authorization"`
}
