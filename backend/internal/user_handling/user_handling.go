package user_handling
import (
	"budgettracker/internal/sql_logic"
	"budgettracker/internal/model"
	"fmt"
    _ "github.com/go-sql-driver/mysql"
	"database/sql"
	"time"

)

//Adds user to database
func Add_user(username string, email string, google_id string, picture_url string){
	res, _ := Find_user(google_id);
	if res == -1 { return }
	if res == 1 {
		fmt.Println("User exists already")
		return
	}
	sts, db := sql_logic.Connect_to_sql();
	if !sts {
		fmt.Println("Error connecting to db-users")
		return
	}
	defer db.Close()

	l_login := time.Now().Format("2006-01-02")
	
	insert_query := "INSERT INTO users (username,email,google_id,picture_url,last_login) VALUES (?, ?, ?, ?, ?)"

	_, err := db.Exec(insert_query, username, email, google_id, picture_url, l_login)
	if err != nil {
				fmt.Println("Couldnt add user:", err)
	}
}

//Returs user
func Find_user(googleID string) (int, model.User){

	res, db := sql_logic.Connect_to_sql();
	if !res {
		fmt.Println("Error connecting to db")
		return -1, model.User{}
	}
	defer db.Close()

	query := `SELECT * FROM users WHERE google_id = ?`

	var user model.User

	err := db.QueryRow(query, googleID).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
		&user.CreatedAt,
        &user.GoogleID,
        &user.PictureURL,
        &user.LastLogin,
    )

	if err == sql.ErrNoRows {
        return 0, model.User{}
    } else if err != nil {
        fmt.Println("Could not perform search for user:", err)
        return -1, model.User{}
    }

    return 1, user
}