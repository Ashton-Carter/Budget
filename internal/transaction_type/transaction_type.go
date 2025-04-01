package transaction_type

import (
	"fmt"
	"budgettracker/internal/model"
	"context"
    "os"
	"bufio"
	"strings"
	"strconv"
	"encoding/csv"

    openai "github.com/sashabaranov/go-openai"
	//GET RID OF THIS BEFORE PRODUCTION
	"github.com/joho/godotenv"
)

type Transaction_type struct {
	Transaction model.Transaction
	T_type string
}

func Get_types(transactions []model.Transaction) []Transaction_type{
	//var with_types []Transaction_type
	var prompt = Create_prompt(transactions)
	var chat_response = AiCall(prompt)
	return Add_types(chat_response, transactions)
}

func Create_prompt(transactions []model.Transaction) string {
	var prompt string = ""
	for i, transaction := range transactions {
		prompt = fmt.Sprintf("%s %d: %s \n", prompt, i, transaction.Description)
	}
	return prompt
}




func AiCall(prompt string) string {
	//GET RID OF THIS BEFORE PRODUCTION
	err := godotenv.Load()
    if err != nil {
        return "Error loading .env file"
    }
	//GET RID OF THIS BEFORE PRODUCTION

    apiKey := os.Getenv("OPENAI_API_KEY")
    client := openai.NewClient(apiKey)
    ctx := context.Background()

    systemMsg := "For each of the following transaction descriptions, give an overall category they fit in.  The categories to choose from must be from the following:Food, Gas, Entertainment, Shopping, Subscriptions, Transfers, Alcohol, Income, Other. Format the output by putting the index number provided, a colon, and then the category."

    req := openai.ChatCompletionRequest{
        Model: "o3-mini", 
        Messages: []openai.ChatCompletionMessage{
            {
                Role:    openai.ChatMessageRoleSystem,
                Content: systemMsg,
            },
            {
                Role:    openai.ChatMessageRoleUser,
                Content: prompt,
            },
        },
    }

    resp, err := client.CreateChatCompletion(ctx, req)
    if err != nil {
        panic(err)
    }

    return resp.Choices[0].Message.Content
}

func Add_types(chat_response string, transactions []model.Transaction) []Transaction_type{
	var transaction_final []Transaction_type

	scanner := bufio.NewScanner(strings.NewReader(chat_response))
	for scanner.Scan() {
		line := scanner.Text()
		int_str := ""
		var text_start int = 0
		var curr_transaction Transaction_type
		for i, letter := range line {
			if letter == ':' {
				text_start = i
				break
			}
			int_str += string(letter)
		}
		index, err := strconv.Atoi(int_str)
		if err != nil {
    		fmt.Println("Error on index:", int_str, "\nError:", err)
    		continue
    	}

		curr_transaction.Transaction = transactions[index]
		curr_transaction.T_type = strings.TrimSpace(line[(text_start + 1):])

		transaction_final = append(transaction_final, curr_transaction)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
	return transaction_final
}


func SaveToCSV(filename string, transactions []Transaction_type) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    
    writer.Write([]string{"Details", "Posting_date", "Description", "Amount", "Type_", "Balance", "Check_Slip"})

    for i, tx := range transactions {
        writer.Write([]string{
            strconv.Itoa(i),
			tx.Transaction.Details,
            tx.Transaction.Posting_date,
			tx.Transaction.Description,
			strconv.FormatFloat(tx.Transaction.Amount, 'f', 2, 64),
			tx.Transaction.Type_,
			strconv.FormatFloat(tx.Transaction.Balance, 'f', 2, 64),
			tx.Transaction.Check_Slip,
            tx.T_type,
        })
    }

    return nil
}