package validations

import (
	"SanjaySinghRajpoot/ticketBooking/config"
	"fmt"
)

func IsUniqueValue(tableName, fieldName, value string) bool {
	var count int64

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = $1", tableName, fieldName)

	row := config.DB.QueryRow(query, value)
	if err := row.Scan(&count); err != nil {
		fmt.Println("Error:", err)
		return false
	}

	return count != 0
}
