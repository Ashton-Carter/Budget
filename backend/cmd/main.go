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
    r.POST("/upload", router_commands.FromCSV)

    //lets add a get for the past 3 months data, maybe just a date range or something

    r.Run(":8080")

    fmt.Println("Done")
}