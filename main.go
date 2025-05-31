package main

import (
	"expense-tracker-cli/cmd"
	"expense-tracker-cli/expense"
	"expense-tracker-cli/storage"
	"log"
	"os"
)

func main() {
	storageManager := storage.NewStorageManager[expense.Expense]("expenses.json")
	expenseManager := expense.NewExpenseManager(storageManager)
	err := cmd.NewCmd(os.Args, expenseManager).Run()
	if err != nil {
		log.Println(err)
	}
}
