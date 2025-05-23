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
	"encoding/csv"
	"strconv"
)

//Sign in logic, more in user_handling.go
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

//Handles getting user
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

//Allows filtering of transactions, sends all transactions to frontend
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

	trans_queries := `
		SELECT 
			t.id,
			t.description, 
			t.posting_date, 
			t.amount, 
			t.type, 
			t.balance, 
			c.name AS category_name
		FROM transactions t
		JOIN categories c ON t.category = c.category_id
		WHERE t.user_id in (select id from users where google_id = ?)
		` 
		
	args := []interface{}{google_id}

	if start_date != "" && end_date != "" {
		trans_queries += ` AND t.posting_date BETWEEN ? AND ?`
		args = append(args, start_date, end_date)
	}
	trans_queries += ` ORDER BY t.posting_date DESC`
	rows, rerr := db.Query(trans_queries, args...)
	var transactions []model.Transaction_type

	if rerr == sql.ErrNoRows {
		fmt.Println("Done, no rows!")
		c.JSON(http.StatusOK, transactions)
		return
	} else if rerr != nil {
		fmt.Println("Database error finding user_id or dates\n", rerr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding user_id or dates"})
		return
	}


	for rows.Next() {
		var tx model.Transaction_type
		if err := rows.Scan(
			&tx.Transaction.Id,
			&tx.Transaction.Description,
			&tx.Transaction.Posting_date,
			&tx.Transaction.Amount,
			&tx.Transaction.Type_,
			&tx.Transaction.Balance,
			&tx.T_type,
		); err != nil {
			fmt.Println("Failed to scan row\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		transactions = append(transactions, tx)
	}
	c.JSON(http.StatusOK, transactions)
}

//Parses csv into transactions
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

//Returns budgets for given user
func GetBudgets(c *gin.Context) {
	google_id := c.Param("google_id")
	connect, db := sql_logic.Connect_to_sql()
	if !connect {
		fmt.Println("Database connection error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	budget_queries := `
		SELECT budgets.budget_id, budgets.user_id, name, budgets.created_at, budget_category.budget_id, category_id, amount
		FROM budgets
		INNER JOIN budget_category
			ON budgets.budget_id = budget_category.budget_id
		INNER JOIN users
			ON budgets.user_id = users.id
		WHERE google_id = ?
		`
	rows, rerr := db.Query(budget_queries, google_id)
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

//Adds budget
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


//Retreives goals for a user from google_id
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

//Adds money to a goal
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

//Adds a goal target
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

//Gets totals for each category
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

//Delete goal
func DeleteGoal(c *gin.Context) {
	goal_id := c.Query("goal_id")
	google_id := c.Query("user_id")
	sql_command := `
	DELETE from goals
	WHERE user_id in (SELECT users.id FROM users WHERE google_id = ?)
	AND goal_id = ?;`
	fmt.Println(goal_id, google_id)


	connect, db := sql_logic.Connect_to_sql()
	if !connect {
		fmt.Println("Database connection error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	_, err := db.Exec(sql_command, google_id, goal_id)
	
	if err != nil {
		fmt.Println("Database error:\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting goal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}


//Deletes a budget
func DeleteBudget(c *gin.Context) {
	budget_id := c.Query("budget_id")
	google_id := c.Query("user_id")
	sql_command := `
	DELETE budgets, budget_category FROM budget_category
	INNER JOIN budgets ON budgets.budget_id = budget_category.budget_id
	WHERE budgets.budget_id = ? and user_id IN (SELECT id FROM users WHERE google_id = ?);`
	fmt.Println(budget_id, google_id)


	connect, db := sql_logic.Connect_to_sql()
	if !connect {
		fmt.Println("Database connection error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	resp, err := db.Exec(sql_command, budget_id, google_id)
	fmt.Println(resp)
	if err != nil {
		fmt.Println("Database error:\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Deleting Budget"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

//Gets all categories/id
func GetCategories(c *gin.Context) {
	sql_command := `SELECT name, category_id from categories;`

	connect, db := sql_logic.Connect_to_sql()
	if !connect {
		fmt.Println("Database connection error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()
	rows, rerr := db.Query(sql_command)
	var cats []model.Category

	if rerr != nil {
		fmt.Println("Database error finding categories", rerr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding categories"})
		return
	}


	for rows.Next() {
		var tx model.Category
		if err := rows.Scan(
			&tx.Name,
			&tx.Category_id,
		); err != nil {
			fmt.Println("Failed to scan row", "\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		cats = append(cats, tx)
	}
	c.JSON(http.StatusOK, cats)

}

//Edits transaction and updates date
func EditTransaction(c *gin.Context) {
	var input model.EditTransactionInput
	transaction_id := c.Param("transaction_id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	
	sql_command := `
	UPDATE transactions
	SET posting_date = ?, description = ?, amount = ?, category = ?, created_at = NOW()
	WHERE user_id IN (SELECT id FROM users WHERE google_id = ?) 
	AND id = ?;`


	connect, db := sql_logic.Connect_to_sql()
	if !connect {
		fmt.Println("Database connection error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	resp, err := db.Exec(sql_command, input.Date, input.Description, input.Amount, input.CategoryID, input.GoogleID, transaction_id)
	fmt.Println(resp)
	if err != nil {
		fmt.Println("Database error:\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Editing transaction"})
		return
	}
	rowsAffected, _ := resp.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No matching transaction found or no change made"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})

}

//Deletes a transaction
func DeleteTransaction(c *gin.Context) {
	transaction_id := c.Param("transaction_id")
	google_id := c.Query("google_id")
	sql_command := `DELETE from transactions
	WHERE id = ?
	AND user_id IN (SELECT id from users WHERE google_id = ?)`

	connect, db := sql_logic.Connect_to_sql()
	if !connect {
		fmt.Println("Database connection error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	resp, err := db.Exec(sql_command, transaction_id, google_id)
	
	if err != nil {
		fmt.Println("Database error:\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Editing transaction"})
		return
	}
	rowsAffected, _ := resp.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No matching transaction found or no change made"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

//Adds a manual entry transaction
func AddTransaction(c *gin.Context) {
	var input model.NewTransactionInput;
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	sql_command := `
	INSERT INTO transactions(user_id, posting_date, description, amount, category, created_at, details, type, balance)
	VALUES((SELECT id from USERS WHERE google_id = ? LIMIT 1), ?, ?, ?, ?, NOW(), '', 'usr_inp', 0.00);
`

	connect, db := sql_logic.Connect_to_sql()
	if !connect {
		fmt.Println("Database connection error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	_, err := db.Exec(sql_command, input.GoogleID, input.Date, input.Description, input.Amount, input.CategoryID)
	
	if err != nil {
		fmt.Println("Database error:\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Editing transaction"})
		return
	}
}

//Save to file
func SaveToFile(c *gin.Context){
	googleID := c.Query("google_id")

	connect, db := sql_logic.Connect_to_sql()
	if !connect {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	sql_command := `
	SELECT Date, description, amount, category from transactions_with_cats
	WHERE GoogleID = ?;`

	rows, err := db.Query(sql_command, googleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
		return
	}
	defer rows.Close()

	// 3. Set headers
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=transactions.csv")
	c.Header("Content-Type", "text/csv")

	// 4. Write CSV
	writer := csv.NewWriter(c.Writer)
	writer.Write([]string{"Date", "Description", "Amount", "Category"})

	for rows.Next() {
		var date, description, category string
		var amount float64
		if err := rows.Scan(&date, &description, &amount, &category); err != nil {
			fmt.Println("Row scan error:", err)
			continue
		}
		writer.Write([]string{date, description, strconv.FormatFloat(amount, 'f', 2, 64), category})
	}
	writer.Flush()
}