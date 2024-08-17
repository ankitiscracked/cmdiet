package diet

import "cmdiet/meals"

type Diet struct {
	MealId         int64
	MealType       string
	Source         string
	Timestamp      int64
	LoggedForDay   int
	LoggedForMonth int
	LoggedForYear  int
}

type DietResp struct {
	Day      string
	MealType string
	Meal     meals.Meal
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
