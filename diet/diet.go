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
	GetLastWeekDiet() []DayDiet
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
	mealId, err := addMeal(d.db, mealName, calories)
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

func addMeal(db *sql.DB, name string, calories int) (int64, error) {
	insertMealSql := `insert into meals (name, calories, timestamp) values (?, ?, ?)`
	statement, err := db.Prepare(insertMealSql)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	result, err := statement.Exec(name, calories, time.Now().UnixMilli())
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.LastInsertId()
}

func (d *dietService) GetLastWeekDiet() []DayDiet {
	rows, err := d.db.Query(`SELECT meal_id, meal_type, timestamp FROM diet WHERE timestamp > ?`, time.Now().AddDate(0, 0, -7).UnixMilli())
	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	type Diet struct {
		name     string
		mealType string
		calories int
	}

	dietMap := make(map[string][]Diet)

	for rows.Next() {
		var mealType string
		var timestamp int64
		var mealId int
		rows.Scan(&mealId, &mealType, &timestamp)
		meal := fetchMealData(mealId)

		day := time.UnixMilli(timestamp).Format("2006-01-02")
		dietMap[day] = append(dietMap[day], Diet{meal.Name, mealType, meal.Calories})
	}

	weekDiets := make([]DayDiet, len(dietMap))

	for day, diets := range dietMap {
		var d DayDiet
		var totalCalories int

		d.Day = day
		for _, diet := range diets {
			switch diet.mealType {
			case "breakfast":
				d.Breakfast = diet.name
			case "lunch":
				d.Lunch = diet.name
			case "dinner":
				d.Dinner = diet.name
			}

			totalCalories += diet.calories
		}
		d.TotalCalories = totalCalories

		weekDiets = append(weekDiets, d)
	}

	return weekDiets
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
