package chase_parser

import (
	"fmt"
	"os"
	"encoding/csv"
    "strconv"

)

type Transaction struct {
	Details string
	Posting_date string
	Description string
	Amount float64
	Type_ string
	Balance float64
	Check_Slip string
}

func ParseCSV(path string) {
	ReadFile(path)
}

func ReadFile(filepath string) {
	file, err := os.Open(filepath)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 //flexible number of columns
	reader.LazyQuotes = true //allows quotes to have commas

    records, err := reader.ReadAll()
    if err != nil {
        fmt.Println("Error reading CSV:", err)
        return
    }

	expectedFields := 7

	for i, record := range records {
		if(i == 0) {
			continue
		}

		if len(record) < expectedFields {
			fmt.Printf("Skipping short row %d\n", i+1)
			continue
		} else if len(record) > expectedFields {
			fmt.Printf("Row %d has extra fields, trimming...\n", i+1)
			record = record[:expectedFields]
		}

		curr_transaction := ReadToTransaction(record)
		printTransaction(curr_transaction);
		break
	}

}

func ReadToTransaction(record []string) Transaction {
	var transaction Transaction
	transaction.Details = record[0]
	transaction.Posting_date = record[1]
	transaction.Description = record[2]

	amount, err := strconv.ParseFloat(record[3], 64)
	if err != nil {
		fmt.Println("Error reading amount")
		amount = 0
	}
	transaction.Amount = amount

	transaction.Type_ = record[4]
	
	balance, err := strconv.ParseFloat(record[5], 64)
	if err != nil {
		fmt.Println("Error reading balance")
		balance = 0
	}

	transaction.Balance = balance
	transaction.Amount = amount
	transaction.Check_Slip = record[6]
	return transaction
}


func printTransaction(trans Transaction) {
	fmt.Println("Details:", trans.Details)
	fmt.Println("Posting Date:", trans.Posting_date) 
	fmt.Println("Transaction:", trans.Description)
	fmt.Println("Amount:", trans.Amount) 
	fmt.Println("Type:", trans.Type_) 
	fmt.Println("Balance:", trans.Balance) 
	fmt.Println("Slip:", trans.Check_Slip)
}