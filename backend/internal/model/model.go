package model

import (
	
)

//Holds all the data types used, Go requires a datatype for sending/recieving data over http
type Transaction struct {
	Id int
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

type Budget struct {
	Budget_id int `json:"budget_id"`
	User_id int `json:"user_id"`
	Name string `json:"name"`
	Created_at string `json:"created_at"`
	Category_id int `json:"category_id"`
	Amount float64 `json:"amount"`
}

type BudgetInput struct{
	User_id string `json:"user_id"`
	Name string `json:"name"`
	Category_id int `json:"category_id"`
	Amount float64 `json:"amount"`
}

type Goal struct {
	Goal_id int `json:"goal_id"`
	User_id int `json:"user_id"`
	Name string `json:"name"`
	Created_at string `json:"created_at"`
	Amount float64 `json:"amount"`
	Current_amount float64 `json:"current_amount"`
}

type GoalInput struct {
	Goal_id int `json:"goal_id"`
	Amount float64 `json:"amount"`
}

type NewGoal struct {
	User_id string `json:"user_id"`
	Name string `json:"name"`
	Amount float64 `json:"amount"`
}

type CategoryTotals struct {
	Name string `json:"name"`
	Total float64 `json:"total"`
}

type Category struct {
	Name string `json:"name"`
	Category_id float64 `json:"category_id"`
}

type EditTransactionInput struct {
	GoogleID   string  `json:"google_id"`
	Date       string  `json:"date"`
	Description string `json:"description"`
	Amount     float64 `json:"amount"`
	CategoryID int     `json:"category_id"`
}

type NewTransactionInput struct {
	GoogleID    string  `json:"google_id"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	CategoryID  int     `json:"category_id"`
}
