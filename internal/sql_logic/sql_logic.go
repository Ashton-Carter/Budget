package sql_logic

import(
	"fmt"
	"os"
	"database/sql"
	"encoding/csv"
    _ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"time"

)

func CSVtoDV(filename string, id int)  {

	res, db := Connect_to_sql()
	if !res {
		fmt.Printf("impossible to create the connection-csv")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf("Couldn't ping database: %s", err)
	}

    f, err := os.Open(filename)
    if err != nil {
        fmt.Println("Unable to read input file " + filename, err)
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
    records, err := csvReader.ReadAll()
    if err != nil {
        fmt.Println("Unable to parse file as CSV for " + filename, err)
    }

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`
	err = db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		fmt.Println("Could not perform search for user:", err)
		return
	}

	if !exists {
		fmt.Println("Couldnt find user")
		return
	}

	


	insert_query := "INSERT INTO transactions (user_id,details,posting_date,description,amount,type,balance,category) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	for index, record := range records{
		parsed, Perr := time.Parse("01/02/2006", record[2])
		if Perr != nil {
			fmt.Println("Error parsing date for sql:", Perr)
		}
		mysqlDate := parsed.Format("2006-01-02")


		_, err := db.Exec(insert_query, id, record[1], mysqlDate, record[3], record[4], record[5], record[6], record[8])
		if err != nil {
			fmt.Println(index, ":", err)
		}
	}
}

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



