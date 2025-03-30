package main

import (
    "budgettracker/internal/csv_parser/chase_parser"
)

func main() {
    chase_parser.ParseCSV("internal/csv_parser/chase_parser/short_test.CSV")
}