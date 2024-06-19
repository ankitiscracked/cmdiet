package diet

import (
	"cmdiet/database"
	"cmdiet/types"
	"fmt"
)

func UpsertMealMacros(mealId int, protien, carbs, fat int) {
	// now := time.Now().Unix()
	// updateSql := `update meals set protein = ?, carbs = ?, fat = ? where id = ?`
	// statement, err := database.DB.Prepare(updateSql)
	//
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("Meal macros has been updated")
}

func GetAllMeals() []types.Meal {
	rows, err := database.DB.Query("select * from meals")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer rows.Close()

	var allMeals []types.Meal

	for rows.Next() {
		var meal types.Meal
		err := rows.Scan(&meal.Id, &meal.Name, &meal.Calories, &meal.Protein, &meal.Carbs, &meal.Fat, &meal.Timestamp)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		allMeals = append(allMeals, meal)
	}

	return allMeals
}
