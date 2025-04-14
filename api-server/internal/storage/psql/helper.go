package psql

import (
	"database/sql"
	"fmt"
)

func zeroRows(res sql.Result, empty error) error {
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't get rows: %w", err)
	}
	if rows == 0 {
		return empty
	}
	return nil
}
