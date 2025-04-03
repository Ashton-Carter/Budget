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
	ID int
	Username string
	Email string
	GoogleID float64
	PictureURL string
	CreatedAt float64
	LastLogin string
}
