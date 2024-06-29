package database

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

func PopulateFakeDiets() {
	insertDietSql := `insert into diet (meal_id, meal_type, source, timestamp) values (?, ?, ?, ?)`
	statement, err := DB.Prepare(insertDietSql)
	if err != nil {
		fmt.Println(err)
	}

	mealTypes := []string{"breakfast", "lunch", "dinner"}
	now := time.Now()

	for i := range 180 {
		date := now.AddDate(0, 0, -i)
		mealId := gofakeit.Number(0, getMealsCout()-1)

		for j := range 2 {
			var timestamp int64

			switch j {
			case 0:
				timestamp = getHourTimestamp(date, gofakeit.Number(6, 12))
			case 1:
				timestamp = getHourTimestamp(date, gofakeit.Number(12, 18))
			case 2:
				timestamp = getHourTimestamp(date, gofakeit.Number(18, 24))
			}

			_, err = statement.Exec(mealId, mealTypes[j], gofakeit.RandomString([]string{"home", "out"}), timestamp)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
	fmt.Println("finished populating diets")
}

func PopulateFakeMeals() {
	insertMealSql := `insert into meals (name, calories, protein, carbs, fats, timestamp) values (?, ?, ?, ?, ?, ?)`
	statement, err := DB.Prepare(insertMealSql)
	if err != nil {
		fmt.Println(err)
	}

	for range 100 {
		name := gofakeit.Breakfast()
		calories := gofakeit.Number(200, 500)
		protein := gofakeit.Number(1, 100)
		carbs := gofakeit.Number(1, 100)
		fats := gofakeit.Number(1, 100)
		_, err := statement.Exec(name, calories, protein, carbs, fats, time.Now().UnixMilli())
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("finished populating meals")
}

func getMealsCout() int {
	var count int
	row := DB.QueryRow("select count(*) from meals")
	err := row.Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	return count
}

func getHourTimestamp(date time.Time, hour int) int64 {
	return time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, date.Location()).UnixMilli()
}
