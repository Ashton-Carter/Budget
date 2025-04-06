package model

import (
	
)

type Transaction struct {
	Details string
	Posting_date string
	Description string
	Amount float64
	Type_ string
	Balance float64
	Check_Slip string
}

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	GoogleID   string `json:"google_id"`
	PictureURL string `json:"picture_url"`
	CreatedAt  string `json:"created_at"`
	LastLogin  string `json:"last_login"`
}

type Transaction_type struct {
	Transaction Transaction
	T_type string
}