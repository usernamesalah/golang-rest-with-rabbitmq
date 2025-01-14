package derrors

import "database/sql"

func HandleSQLError(err error, format string, args ...any) error {
	if err != nil {
		if err == sql.ErrNoRows {
			return New(NotFound, "Not found")
		}
		return WrapStack(err, Unknown, format, args...)
	}
	return nil
}
