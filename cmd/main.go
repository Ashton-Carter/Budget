package main

import (
    // "budgettracker/internal/csv_parser/chase_parser"
    // "budgettracker/internal/transaction_type"
    // "budgettracker/internal/sql_logic"
    //"budgettracker/internal/user_handling"
    "budgettracker/internal/router_commands"
    "fmt"
    "github.com/gin-gonic/gin"

)

func main() {
    // transactions := chase_parser.ParseCSV("internal/csv_parser/test.CSV")
    // type_trans := transaction_type.Get_types(transactions)
    // transaction_type.SaveToCSV("dtb.csv", type_trans)
    // sql_logic.CSVtoDV("dtb.csv", 1)
    //user_handling.Add_user("GodEmperor", "NO", "BEAST", "GOD_EMPEROR.png")
    r := gin.Default()
    r.POST("/auth/google", router_commands.GoogleAuth)
    r.GET("/users/:google_id", router_commands.GetUser)
    r.GET("/transactions/:google_id", router_commands.GetTransactions)

    r.Run(":8080")

    fmt.Println("Done")
}