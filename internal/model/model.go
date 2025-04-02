package model

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
	id int
	username string
	email string
	google_id float64
	picture_url string
	created_at float64
	last_login string
}