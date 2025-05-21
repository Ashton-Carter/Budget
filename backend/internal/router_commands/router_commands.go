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
		var ignored string;
		if err := rows.Scan(
			&bud.Budget_id,
			&bud.User_id,
			&bud.Name,
			&bud.Created_at,
			&ignored,
			&bud.Category_id,
			&bud.Amount,
		); err != nil {
			fmt.Println("Failed to scan row:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		budgets = append(budgets, bud)
	}
	c.JSON(http.StatusOK, budgets)
}

func AddBudget(c *gin.Context) {
	var input model.BudgetInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	connect, db := sql_logic.Connect_to_sql()
	if !connect {
		fmt.Println("Database connection error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	var q_id string
	usr_query := "Select id from users WHERE google_id = ?"
	err := db.QueryRow(usr_query, input.User_id).Scan(&q_id)

	if err == sql.ErrNoRows {
		fmt.Println("User not found in users")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found in users"})
		return
	} else if err != nil {
		fmt.Println("Database error finding user_id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding user_id"})
		return
	}

	var budgetExists bool

	budget_query := `SELECT budget_id FROM budgets
WHERE user_id = ?
AND name = ? `
	var bud_id int
	err = db.QueryRow(budget_query, q_id, input.Name).Scan(&bud_id)

	if err == sql.ErrNoRows {
		budgetExists = false
	} else if err != nil {
		fmt.Println("Database error finding budget")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	} else {
		budgetExists = true
	}

	if !budgetExists {
		var budget_input = `INSERT INTO budgets(user_id, name, created_at) VALUES (?, ?, NOW())`
		_, err = db.Exec(budget_input, q_id, input.Name)
	
		if err != nil {
			fmt.Println("Database error inserting budget")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		err = db.QueryRow(budget_query, q_id, input.Name).Scan(&bud_id)
		if err != nil {
			fmt.Println("Database error finding budget after insert")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding budget"})
			return
		}
	}

	
	var insert_query = `INSERT INTO budget_category(budget_id, category_id, amount) VALUES (?, ?, ?)`
	_, err = db.Exec(insert_query, bud_id, input.Category_id, input.Amount)
	
	if err != nil {
		fmt.Println("Database error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}




	c.JSON(http.StatusOK, gin.H{
		"message": "Budget created successfully",
		"data":    input,
	})
}

func GetGoals(c *gin.Context) {
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

	goal_query := `
		SELECT *
		FROM goals
		WHERE user_id = ?
		`
	rows, rerr := db.Query(goal_query, q_id)
	var goals []model.Goal

	if rerr == sql.ErrNoRows {
		fmt.Println("Done, no rows!")
		c.JSON(http.StatusOK, goals)
		return
	} else if rerr != nil {
		fmt.Println("Database error finding user_id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding user_id"})
		return
	}


	for rows.Next() {
		var goal model.Goal
		if err := rows.Scan(
			&goal.Goal_id,
			&goal.User_id,
			&goal.Name,
			&goal.Created_at,
			&goal.Amount,
			&goal.Current_amount,
		); err != nil {
			fmt.Println("Failed to scan row:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		goals = append(goals, goal)
	}
	c.JSON(http.StatusOK, goals)
}

func AddToGoal(c *gin.Context) {
	var req model.GoalInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}


	connect, db := sql_logic.Connect_to_sql()
	if !connect {
		fmt.Println("Database connection error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	goal_query := `
		UPDATE goals
		SET current_amount = current_amount + ?
		WHERE goal_id = ?
		`
	_, err := db.Exec(goal_query, req.Amount, req.Goal_id)
	

	if err != nil {
		fmt.Println("Database error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, "Submitted")
	
}

func AddGoal(c *gin.Context){
	

	var input model.NewGoal
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	google_id := input.User_id
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

	inp_query :=`INSERT INTO goals(user_id, name, created_at, amount, current_amount)
				VALUES(?, ?, NOW(), ?, 0)`

	_, err = db.Exec(inp_query, q_id, input.Name, input.Amount)
	
	if err != nil {
		fmt.Println("Database error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating goal"})
		return
	}

}

func Get_Monthly_Totals(c *gin.Context){
	sql_command := 
	`SELECT categories.name, sum(transactions.amount) as total
    from users
    INNER JOIN transactions
             ON users.id = transactions.user_id
    INNER JOIN categories
             ON transactions.category = categories.category_id
    WHERE google_id = ?
    AND transactions.posting_date BETWEEN ? AND ?
	AND transactions.amount < 0
    GROUP BY categories.name;`


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
	fmt.Println(google_id, start_date, end_date)
	rows, rerr := db.Query(sql_command, google_id, start_date, end_date)
	var cat_totals []model.CategoryTotals

	if rerr == sql.ErrNoRows {
		fmt.Println("Done, no rows!")
		c.JSON(http.StatusOK, cat_totals)
		return
	} else if rerr != nil {
		fmt.Println("Database error finding user_id or dates", rerr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding user_id or dates"})
		return
	}


	for rows.Next() {
		var tx model.CategoryTotals
		if err := rows.Scan(
			&tx.Name,
			&tx.Total,
		); err != nil {
			fmt.Println("Failed to scan row", "\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		cat_totals = append(cat_totals, tx)
	}
	c.JSON(http.StatusOK, cat_totals)
}