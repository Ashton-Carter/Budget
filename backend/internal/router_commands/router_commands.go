package router_commands

import (
	"github.com/gin-gonic/gin"
	"budgettracker/internal/sql_logic"
	"budgettracker/internal/user_handling"
	"budgettracker/internal/model"
	"budgettracker/internal/csv_parser/chase_parser"
	"budgettracker/internal/transaction_type"
	"net/http"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
	"fmt"
)

func GoogleAuth(c *gin.Context) {
	var input model.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	res, user := user_handling.Find_user(input.GoogleID)
	

	if res ==  0{
		user_handling.Add_user(input.Username, input.Email, input.GoogleID, input.PictureURL)
		res, user = user_handling.Find_user(input.GoogleID)
		if res != 1 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add User..."})
			return
		}
	} else if res == -1 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUser(c *gin.Context) {
	google_id := c.Param("google_id")

	res, user := user_handling.Find_user(google_id)
	

	if res ==  0{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User Not Found"})
		return
	} else if res == -1 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	sql_logic.LastLoginUpdate(google_id)

	c.JSON(http.StatusOK, user)
	
}

func GetTransactions(c *gin.Context) {
	google_id := c.Param("google_id")
	start_date := c.Query("start_date")
	end_date := c.Query("end_date")
	connect, db := sql_logic.Connect_to_sql()
	if !connect {
		fmt.Println("Database connection error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	var q_id string
	usr_query := "Select id from users WHERE google_id = ?"
	err := db.QueryRow(usr_query, google_id).Scan(&q_id)

	if err == sql.ErrNoRows {
		fmt.Println("User not found in users")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found in users"})
		return
	} else if err != nil {
		fmt.Println("Database error finding user_id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding user_id"})
		return
	}

	trans_queries := `
		SELECT 
			t.details, 
			t.posting_date, 
			t.amount, 
			t.type, 
			t.balance, 
			c.name AS category_name
		FROM transactions t
		JOIN categories c ON t.category = c.category_id
		WHERE t.user_id = ? 
		AND t.posting_date BETWEEN ? AND ?
		`
	rows, rerr := db.Query(trans_queries, q_id, start_date, end_date)
	var transactions []model.Transaction_type

	if rerr == sql.ErrNoRows {
		fmt.Println("Done, no rows!")
		c.JSON(http.StatusOK, transactions)
		return
	} else if rerr != nil {
		fmt.Println("Database error finding user_id or dates")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding user_id or dates"})
		return
	}


	for rows.Next() {
		var tx model.Transaction_type
		if err := rows.Scan(
			&tx.Transaction.Details,
			&tx.Transaction.Posting_date,
			&tx.Transaction.Amount,
			&tx.Transaction.Type_,
			&tx.Transaction.Balance,
			&tx.T_type,
		); err != nil {
			fmt.Println("Failed to scan row")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		transactions = append(transactions, tx)
	}
	c.JSON(http.StatusOK, transactions)
}

func FromCSV(c *gin.Context) {
    googleID := c.PostForm("google_id")
    fmt.Println("Google ID:", googleID)

    uploadedFile, err := c.FormFile("file")
    if err != nil {
		fmt.Println("No file uploaded")
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
        return
    }

    file, err := uploadedFile.Open()
    if err != nil {
		fmt.Println("Unable to open file")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open file"})
        return
    }
    defer file.Close()

	transactions := chase_parser.ParseCSVFile(file)
	fmt.Println("Parsed successfully")
    type_trans := transaction_type.Get_types(transactions)
	fmt.Println("Got types")
    sql_logic.TranstoDV(type_trans, googleID)
	fmt.Println("All done!")

	c.JSON(http.StatusOK, gin.H{"message": "Success"})

}

func GetBudgets(c *gin.Context) {
	google_id := c.Param("google_id")
	connect, db := sql_logic.Connect_to_sql()
	if !connect {
		fmt.Println("Database connection error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	var q_id string
	usr_query := "Select id from users WHERE google_id = ?"
	err := db.QueryRow(usr_query, google_id).Scan(&q_id)

	if err == sql.ErrNoRows {
		fmt.Println("User not found in users")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found in users"})
		return
	} else if err != nil {
		fmt.Println("Database error finding user_id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding user_id"})
		return
	}

	budget_queries := `
		SELECT *
		FROM budgets
		INNER JOIN budget_category
			ON budgets.budget_id = budget_category.budget_id
		WHERE user_id = ?
		`
	rows, rerr := db.Query(budget_queries, q_id)
	var budgets []model.Budget

	if rerr == sql.ErrNoRows {
		fmt.Println("Done, no rows!")
		c.JSON(http.StatusOK, budgets)
		return
	} else if rerr != nil {
		fmt.Println("Database error finding user_id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding user_id"})
		return
	}


	for rows.Next() {
		var bud model.Budget
		if err := rows.Scan(
			&bud.Budget_id,
			&bud.User_id,
			&bud.Name,
			&bud.Created_at,
			&bud.Category_id,
			&bud.Amount,
		); err != nil {
			fmt.Println("Failed to scan row")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		budgets = append(budgets, bud)
	}
	c.JSON(http.StatusOK, budgets)
}