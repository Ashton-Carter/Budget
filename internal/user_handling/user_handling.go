package user_handling
import (
	"golang.org/x/crypto/bcrypt"
	"budgettracker/internal/model"
)

func add_user(username string, email string, google_id string, picture_url string, last_login string){
	if !find_user(google_id)[0] {
		fmt.Println("User Exists Already")
		return
	}
}

func find_user(googleID string) (bool, model.User){
	query := `SELECT * FROM users WHERE google_id = ?`

	var user model.User

	err := db.QueryRow(query, googleID).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.GoogleID,
        &user.PictureURL,
        &user.CreatedAt,
        &user.LastLogin,
    )

	if err == sql.ErrNoRows {
        return false, model.User{}
    } else if err != nil {
        fmt.Println("Could not perform search for user:", err)
        return false, model.User{}
    }

    return true, user
}