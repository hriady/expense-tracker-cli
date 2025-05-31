package cmd

import (
	"errors"
	"expense-tracker-cli/expense"
	"flag"
	"fmt"
	"time"
)

type Cmd struct {
	args           []string
	expenseManager expenseManager
}

type expenseManager interface {
	Add(expense expense.Expense) (expense.Expense, error)
	Update(expense expense.Expense) (expense.Expense, error)
	List()
	Delete(id int64)
	Summary(month int64)
}

func NewCmd(args []string, manager expenseManager) *Cmd {
	return &Cmd{
		args:           args,
		expenseManager: manager,
	}
}

func (c *Cmd) Run() error {
	if len(c.args) < 2 {
		return errors.New("missing arguments")
	}
	
	command := c.args[1]
	
	switch command {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		
		description := addCmd.String("description", "", "description")
		amount := addCmd.Int64("amount", 0, "amount")
		err := addCmd.Parse(c.args[2:])
		if err != nil {
			return errors.New("invalid arguments for add command")
		}
		data := expense.Expense{
			Description: *description,
			Amount:      *amount,
			Date:        time.Now(),
		}
		
		if data.Amount <= 0 {
			return errors.New("amount must be greater than 0")
		}
		
		if data.Description == "" {
			return errors.New("description is required")
		}
		
		inserted, err := c.expenseManager.Add(data)
		if err != nil {
			return errors.New("error adding expense")
		}
		fmt.Printf("Expense added successfully (ID:%d)\n", inserted.ID)
	
	case "update":
		updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
		
		id := updateCmd.Int64("id", 0, "id")
		description := updateCmd.String("description", "", "description")
		amount := updateCmd.Int64("amount", 0, "amount")
		date := updateCmd.String("date", "", "date")
		err := updateCmd.Parse(c.args[2:])
		if err != nil {
			return errors.New("invalid arguments for update command")
		}
		
		if id == nil || *id == 0 {
			return errors.New("id is required")
		}
		
		var expenseDate time.Time
		if date != nil && *date != "" {
			parsedDate, err := time.Parse("2006-01-02", *date)
			if err != nil {
				return errors.New("invalid date format. Use YYYY-MM-DD")
			}
			expenseDate = parsedDate
		}
		
		data := expense.Expense{
			ID:          *id,
			Description: *description,
			Amount:      *amount,
			Date:        expenseDate,
		}
		
		if data.Amount <= 0 {
			return errors.New("amount must be greater than 0")
		}
		
		if data.Description == "" {
			return errors.New("description is required")
		}
		
		updated, err := c.expenseManager.Update(data)
		if err != nil {
			return errors.New("error updating expense")
		}
		fmt.Printf("Expense updated successfully (ID:%d, Description: %s, Amount: %d)\n", updated.ID, updated.Description, updated.Amount)
	
	case "list":
		
		c.expenseManager.List()
	
	case "delete":
		
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		id := deleteCmd.Int64("id", 0, "id")
		err := deleteCmd.Parse(c.args[2:])
		if err != nil {
			return errors.New("invalid arguments for delete command")
		}
		
		if id == nil || *id == 0 {
			return errors.New("id is required")
		}
		
		c.expenseManager.Delete(*id)
	
	case "summary":
		
		summaryCmd := flag.NewFlagSet("summary", flag.ExitOnError)
		month := summaryCmd.Int64("month", 0, "month")
		err := summaryCmd.Parse(c.args[2:])
		if err != nil {
			return errors.New("invalid arguments for summary command")
		}
		
		if month == nil || *month == 0 {
			c.expenseManager.Summary(0)
		} else {
			c.expenseManager.Summary(*month)
		}
	
	default:
	}
	
	return nil
}
