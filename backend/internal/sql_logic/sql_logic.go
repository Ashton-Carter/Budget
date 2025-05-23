package sql_logic

import(
	"fmt"
	"os"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"time"
	"budgettracker/internal/model"
)

//Returns connection to sql
func Connect_to_sql() (bool, *sql.DB){
	load := godotenv.Load()
	if load != nil {
		return false, nil
	}

	dsn := os.Getenv("SQL_URL")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return false, nil
	}
	return true, db
}

//Updates last login
func LastLoginUpdate(google_id string){
	sql_command := `
	UPDATE users
SET last_login = NOW()
WHERE google_id = ?;
`

	connect, db := Connect_to_sql()
	if !connect {
		fmt.Println("Database connection error")
		return
	}
	defer db.Close()

	_, err := db.Exec(sql_command, google_id)
	
	if err != nil {
		fmt.Println("Database error:\n", err)
		return
	}
}



//Inserst transaction to database
func TranstoDV(records []model.Transaction_type, googleID string) {
	res, db := Connect_to_sql()
	if !res {
		fmt.Println("impossible to create the connection-csv")
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf("Couldn't ping database: %s", err)
		return
	}

	var userID int
	query := `SELECT id FROM users WHERE google_id = ?`
	err := db.QueryRow(query, googleID).Scan(&userID)
	if err != nil {
		fmt.Println("Could not find user ID from google_id:", err)
		return
	}

	insertQuery := `INSERT INTO transactions (user_id, details, posting_date, description, amount, type, balance, category) 
	                VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	for index, record := range records {
		parsed, parseErr := time.Parse("01/02/2006", record.Transaction.Posting_date)
		if parseErr != nil {
			fmt.Println("Error parsing date for SQL:", parseErr)
			continue
		}
		mysqlDate := parsed.Format("2006-01-02")

		_, err := db.Exec(insertQuery, userID, record.Transaction.Details, mysqlDate, record.Transaction.Description,
			record.Transaction.Amount, record.Transaction.Type_, record.Transaction.Balance, record.T_type)

		if err != nil {
			fmt.Println(index, ":", err)
		}
	}
}
