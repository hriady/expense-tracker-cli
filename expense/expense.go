package expense

import (
	"fmt"
	"time"
)

type ExpenseManager struct {
	storageManager storageManager[Expense]
	expenses       []Expense
}

type storageManager[T any] interface {
	Save(data []T) error
	Load() (resp []T, err error)
}

func NewExpenseManager(manager storageManager[Expense]) *ExpenseManager {
	expenseManager := ExpenseManager{
		storageManager: manager,
	}
	
	expenses, _ := expenseManager.storageManager.Load()
	expenseManager.expenses = expenses
	
	return &expenseManager
}

func (em *ExpenseManager) Add(data Expense) (Expense, error) {
	data.ID = em.findNextID()
	em.expenses = append(em.expenses, data)
	
	err := em.storageManager.Save(em.expenses)
	if err != nil {
		return Expense{}, err
	}
	
	return data, nil
}

func (em *ExpenseManager) Update(data Expense) (Expense, error) {
	
	for i, expense := range em.expenses {
		if expense.ID == data.ID {
			em.expenses[i] = data
			break
		}
	}
	
	err := em.storageManager.Save(em.expenses)
	if err != nil {
		return Expense{}, err
	}
	
	return data, nil
}

func (em *ExpenseManager) List() {
	fmt.Printf("# %-3s %-10s %-12s %s\n", "ID", "Date", "Description", "Amount")
	for _, expense := range em.expenses {
		fmt.Printf("# %-3d %-10s %-12s %d\n", expense.ID, expense.Date.Format("2006-01-02"), expense.Description, expense.Amount)
	}
	fmt.Println()
}

func (em *ExpenseManager) Delete(id int64) {
	newExpenses := []Expense{}
	for _, expense := range em.expenses {
		if expense.ID != id {
			newExpenses = append(newExpenses, expense)
		}
	}
	
	em.expenses = newExpenses
	em.storageManager.Save(em.expenses)
	fmt.Println("Expense deleted successfully")
}

func (em *ExpenseManager) Summary(month int64) {
	total := 0
	monthStr := time.Month(month).String()
	if month > 0 {
		for _, expense := range em.expenses {
			if int64(expense.Date.Month()) == month {
				total += int(expense.Amount)
			}
		}
		fmt.Printf("Total expenses for %s: $%d\n", monthStr, total)
	} else {
		for _, expense := range em.expenses {
			total += int(expense.Amount)
		}
		fmt.Printf("Total expenses: $%d\n", total)
	}
}

func (em *ExpenseManager) findNextID() int64 {
	currentID := int64(0)
	
	for _, expense := range em.expenses {
		if expense.ID > currentID {
			currentID = expense.ID
		}
	}
	
	currentID++
	return currentID
}
