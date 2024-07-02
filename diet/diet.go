package diet

import (
	"cmdiet/database"
	"cmdiet/types"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type DietService interface {
	LogDiet(mealType string, mealName string, calories int, source string) error
	MealTypeLoggedForToday(mealType string) bool
	GetDayDietsByOffset(offset int, timestamp int64) ([]DayDiet, error)
	GetBatchedDayDiets(afterUnixMilli int64, beforeUnixMilli int64) ([]DayDiet, error)
}

type dietService struct {
	db *sql.DB
}

type DayDiet struct {
	Day           string
	Breakfast     string
	Lunch         string
	Dinner        string
	Protein       int
	Carbs         int
	Fat           int
	TotalCalories int
}

func NewDietService(db *sql.DB) DietService {
	return &dietService{db}
}

var DefaultDietService = NewDietService(database.DB)

func (d *dietService) LogDiet(mealType string, mealName string, calories int, source string) error {
	if d.MealTypeLoggedForToday(mealType) {
		return errors.New("meal type already logged for today")
	}
	mealId, err := AddMeal(d.db, mealName, calories)
	if err != nil {
		fmt.Println(err)
	}

	insertDiet(d.db, mealId, mealType, source)
	fmt.Println("Diet logged successfully")
	return nil
}

func insertDiet(db *sql.DB, mealId int64, mealType string, source string) {
	timestamp := time.Now().UnixMilli()
	insertDietSql := `insert into diet (meal_id, meal_type, source, timestamp) values (?, ?, ?, ?)`
	statement, err := db.Prepare(insertDietSql)
	if err != nil {
		fmt.Println(err)
	}

	_, err = statement.Exec(mealId, mealType, source, timestamp)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (d *dietService) GetDayDietsByOffset(offset int, timestamp int64) ([]DayDiet, error) {
	if timestamp != 0 {
		boundingMilli := time.UnixMilli(timestamp).AddDate(0, 0, offset).UnixMilli()
		if offset < 0 {
			return d.GetBatchedDayDiets(boundingMilli, timestamp)
		} else {
			return d.GetBatchedDayDiets(timestamp, boundingMilli)
		}
	} else {
		if offset < 0 {
			return nil, errors.New("offset must be more than 0")
		}
		afterUnixMilli := time.Now().AddDate(0, 0, -offset).UnixMilli()
		return d.GetBatchedDayDiets(afterUnixMilli, 0)
	}
}

func (d *dietService) GetBatchedDayDiets(afterUnixMilli int64, beforeUnixMilli int64) ([]DayDiet, error) {
	rows, err := getBatchedDiets(afterUnixMilli, beforeUnixMilli, d)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer rows.Close()

	dietMap := dietsOfDay(rows)
	var weekDiets []DayDiet

	return weeklyDiets(dietMap, weekDiets), nil
}

func getBatchedDiets(afterUnixMilli int64, beforeUnixMilli int64, d *dietService) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error
	if afterUnixMilli != 0 && beforeUnixMilli != 0 {
		rows, err = d.db.Query(`SELECT meal_id, meal_type, timestamp FROM diet WHERE timestamp > ? and timestamp < ?`, afterUnixMilli, beforeUnixMilli)
	} else if afterUnixMilli != 0 {
		rows, err = d.db.Query(`SELECT meal_id, meal_type, timestamp FROM diet WHERE timestamp > ?`, afterUnixMilli)
	} else if beforeUnixMilli != 0 {
		rows, err = d.db.Query(`SELECT meal_id, meal_type, timestamp FROM diet WHERE timestamp < ?`, beforeUnixMilli)
	}
	return rows, err
}

/*
this function takes a map of day - diets of that day for a week and
returns tabular view of the details of the diets
*/
func weeklyDiets(dietMap map[string][]types.Diet, weekDiets []DayDiet) []DayDiet {
	for day, diets := range dietMap {
		var d DayDiet
		var totalCalories int

		d.Day = day
		for _, diet := range diets {
			d.Protein = int(diet.Meal.Protein.Int16)
			d.Carbs = int(diet.Meal.Carbs.Int16)
			d.Fat = int(diet.Meal.Fat.Int16)
			switch diet.MealType {
			case "breakfast":
				d.Breakfast = diet.Meal.Name
			case "lunch":
				d.Lunch = diet.Meal.Name
			case "dinner":
				d.Dinner = diet.Meal.Name
			}

			totalCalories += diet.Meal.Calories
		}
		d.TotalCalories = totalCalories

		weekDiets = append(weekDiets, d)
	}

	return weekDiets
}

// instead of this, use the group by clause in the query
func dietsOfDay(rows *sql.Rows) map[string][]types.Diet {
	dietMap := make(map[string][]types.Diet)

	for rows.Next() {
		var mealType string
		var timestamp int64
		var mealId int
		rows.Scan(&mealId, &mealType, &timestamp)
		meal := fetchMealData(mealId)

		day := time.UnixMilli(timestamp).Format("2006-01-02")
		dietMap[day] = append(dietMap[day], types.Diet{Day: day, MealType: mealType, Meal: meal})
	}
	return dietMap
}

func (d *dietService) MealTypeLoggedForToday(mealType string) bool {
	var count int
	start, end := timeStampRangeForToday()
	err := d.db.QueryRow(`select count(*) from diet where meal_type = ? and timestamp > ? and timestamp < ?`, mealType, start, end).Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	return count > 0
}

func timeStampRangeForToday() (int64, int64) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	return startOfDay.UnixMilli(), endOfDay.UnixMilli()
}

func fetchMealData(mealId int) types.Meal {
	row := database.DB.QueryRow(`select name, calories, protein, carbs, fats from meals where id = ?`, mealId)
	var meal types.Meal
	err := row.Scan(&meal.Name, &meal.Calories, &meal.Protein, &meal.Carbs, &meal.Fat)
	if err != nil {
		fmt.Println(err)
	}
	return meal
}
