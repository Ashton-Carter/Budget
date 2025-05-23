package main

import (
    "budgettracker/internal/router_commands"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"

)

func main() {
    //All endpoints for the frontend to connect to
    r := gin.Default()
    r.Use(cors.Default())
    r.POST("/auth/google", router_commands.GoogleAuth)
    r.GET("/users/:google_id", router_commands.GetUser)
    r.GET("/transactions/:google_id", router_commands.GetTransactions)
    r.GET("/categorytotals/:google_id", router_commands.Get_Monthly_Totals)
    r.GET("/budgets/:google_id", router_commands.GetBudgets)
    r.GET("/categories", router_commands.GetCategories)
    r.POST("/upload", router_commands.FromCSV)
    r.POST("/budgets", router_commands.AddBudget)
    r.GET("/goals/:google_id", router_commands.GetGoals)
    r.POST("/goals/add", router_commands.AddToGoal)
    r.POST("/goals", router_commands.AddGoal)
    r.DELETE("/goals/", router_commands.DeleteGoal)
    r.DELETE("/budgets/", router_commands.DeleteBudget)
    r.PUT("/transactions/:transaction_id", router_commands.EditTransaction)
    r.DELETE("/transactions/:transaction_id", router_commands.DeleteTransaction)
    r.POST("/transactions", router_commands.AddTransaction)
    r.GET("/transactions/download", router_commands.SaveToFile)
    r.Run(":8080")
    
}