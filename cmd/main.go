package main

import (
    // "budgettracker/internal/csv_parser/chase_parser"
    // "budgettracker/internal/transaction_type"
    // "budgettracker/internal/sql_logic"
    "budgettracker/internal/user_handling"
    "fmt"
)

func main() {
    // transactions := chase_parser.ParseCSV("internal/csv_parser/test.CSV")
    // type_trans := transaction_type.Get_types(transactions)
    // transaction_type.SaveToCSV("dtb.csv", type_trans)
    // sql_logic.CSVtoDV("dtb.csv", 1)
    user_handling.Add_user("GodEmperor", "NO", "BEAST", "GOD_EMPEROR.png")
    fmt.Println("Done")
}