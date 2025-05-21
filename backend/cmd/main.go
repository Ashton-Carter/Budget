package main

import (
    // "budgettracker/internal/csv_parser/chase_parser"
    // "budgettracker/internal/transaction_type"
    // "budgettracker/internal/sql_logic"
    //"budgettracker/internal/user_handling"
    "budgettracker/internal/router_commands"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"

)

func main() {
    r := gin.Default()
    r.Use(cors.Default())
    r.POST("/auth/google", router_commands.GoogleAuth)
    r.GET("/users/:google_id", router_commands.GetUser)
    r.GET("/transactions/:google_id", router_commands.GetTransactions)
    r.GET("/categorytotals/:google_id", router_commands.Get_Monthly_Totals)
    r.GET("/budgets/:google_id", router_commands.GetBudgets)
    r.POST("/upload", router_commands.FromCSV)
    r.POST("/budgets", router_commands.AddBudget)
    r.GET("/goals/:google_id", router_commands.GetGoals)
    r.POST("/goals/add", router_commands.AddToGoal)
    r.POST("/goals", router_commands.AddGoal)

    r.Run(":8080")
    

    fmt.Println("Done")
}