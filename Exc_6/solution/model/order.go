package model

import (
	"fmt"
	"time"
)

const (
	orderFilename = "order_%d.md"

	// Updated markdownTemplate to include Drink Name instead of just Drink ID
	markdownTemplate = `# Order: %d

| Created At      | Drink Name   | Amount |
|-----------------|-------------|--------|
| %s | %s | %d |

Thanks for drinking with us!
`
)

type Order struct {
	Base
	Amount uint64 `json:"amount"`
	// Relationships
	// foreign key
	DrinkID uint  `json:"drink_id" gorm:"not null"`
	Drink   Drink `json:"drink"`
}

// ToMarkdown returns a markdown-formatted receipt for the order.
func (o *Order) ToMarkdown() string {
	drinkName := o.Drink.Name
	if drinkName == "" {
		drinkName = fmt.Sprintf("ID %d", o.DrinkID)
	}
	return fmt.Sprintf(
		markdownTemplate,
		o.ID,
		o.CreatedAt.Format(time.Stamp), // e.g. "Nov 12 17:12:39"
		drinkName,
		o.Amount,
	)
}

func (o *Order) GetFilename() string {
	return fmt.Sprintf(orderFilename, o.ID)
}
