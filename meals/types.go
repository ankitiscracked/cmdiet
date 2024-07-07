package meals

import "database/sql"

type Meal struct {
	Id        int
	Name      string
	Calories  int
	Protein   sql.NullInt16
	Carbs     sql.NullInt16
	Fat       sql.NullInt16
	Timestamp int
}
