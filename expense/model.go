package expense

import "time"

type Expense struct {
	ID          int64     `json:"id"`
	Description string    `json:"description"`
	Amount      int64     `json:"amount"`
	Date        time.Time `json:"date"`
}
