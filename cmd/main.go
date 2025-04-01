package main

import (
    "budgettracker/internal/csv_parser/chase_parser"
    "budgettracker/internal/transaction_type"
    "fmt"
)

func main() {
    transactions := chase_parser.ParseCSV("internal/csv_parser/test.CSV")
    type_trans := transaction_type.Get_types(transactions)
    transaction_type.SaveToCSV("dtb.csv", type_trans)
    fmt.Println("Done")
}