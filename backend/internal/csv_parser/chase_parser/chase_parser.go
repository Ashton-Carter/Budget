package chase_parser

import (
	"fmt"
	"encoding/csv"
    "strconv"
	"budgettracker/internal/model"
	"time"
	"mime/multipart"
)


//Header for the chase CSV file
var correct_types [7]string= [7]string{"Details","Posting Date","Description","Amount","Type","Balance","Check or Slip #"}


//Parses csv row to data type
func ReadToTransaction(record []string) model.Transaction {
	var transaction model.Transaction
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

//Used for debugging
func printTransaction(trans model.Transaction) {
	fmt.Println("Details:", trans.Details)
	fmt.Println("Posting Date:", trans.Posting_date) 
	fmt.Println("Transaction:", trans.Description)
	fmt.Println("Amount:", trans.Amount) 
	fmt.Println("Type:", trans.Type_) 
	fmt.Println("Balance:", trans.Balance) 
	fmt.Println("Slip:", trans.Check_Slip)
}

//Transforms date format, used to limit data size
func WithinLast3Months(dateStr string) bool {
    layout := "01/02/2006" // for MM/DD/YYYY
    postingDate, err := time.Parse(layout, dateStr)
    if err != nil {
        fmt.Println("Error parsing date:", err)
        return false
    }

    now := time.Now()

    startOfRange := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -3, 0)

    endOfRange := now

    return !postingDate.Before(startOfRange) && !postingDate.After(endOfRange)
}


//temporary test from filepath func
func ParseCSVFile(file multipart.File) []model.Transaction{
	var transactions []model.Transaction
	
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 //flexible number of columns
	reader.LazyQuotes = true //allows quotes to have commas

    records, err := reader.ReadAll()
    if err != nil {
        fmt.Println("Error reading CSV:", err)
        return nil
    }

	expectedFields := 7

	for i, record := range records {
		if(i == 0) {
			for i, tpe := range correct_types {
				if tpe != record[i] {
					fmt.Println("Incorrect File Format")
					return nil
				}
			}
			continue
		}

		if len(record) < expectedFields {
			fmt.Printf("Skipping short row %d\n", i+1)
			continue
		} else if len(record) > expectedFields {
			//fmt.Printf("Row %d has extra fields, trimming...\n", i+1)
			record = record[:expectedFields]
		}

		curr_transaction := ReadToTransaction(record)
		if WithinLast3Months(curr_transaction.Posting_date){	
			transactions = append(transactions, curr_transaction)
		}
		// fmt.Println("Added")
	}
	return transactions

}