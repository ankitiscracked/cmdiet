package diet

import (
	"cmdiet/meals"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"gorm.io/gorm"
)

type MealType int

const (
	Breakfast MealType = iota + 1
	Lunch
	Dinner
)

func (m MealType) String() string {
	switch m {
	case Breakfast:
		return "breakfast"
	case Lunch:
		return "lunch"
	case Dinner:
		return "dinner"
	default:
		return fmt.Sprintf("unknown meal type (%d)", int(m))
	}
}

var DS *DietServiceImpl

type DietService interface {
	LogDiet(mealType string, mealName string, calories int, source string) error
	MealTypeLoggedForToday(mealType string) bool
	GetDayDietsByOffset(offset int, timestamp int64) ([]DayDiet, error)
	GetBatchedDayDiets(afterUnixMilli int64, beforeUnixMilli int64) ([]DayDiet, error)
}

type DietServiceImpl struct {
	DB *gorm.DB
}

func (d *DietServiceImpl) LogDietWithNewMeal(mealType MealType, mealName string, source string, logForDate time.Time) error {
	if !ValidMealType(mealType) {
		return fmt.Errorf("invalid meal type %v, please use breakfast, lunch, snacks, or dinner", mealType)
	}

	logged, err := d.MealTypeLoggedForToday(mealType)
	if err != nil {
		return err
	}

	if logged {
		return errors.New("meal type already logged for today")
	}

	meal, err := meals.MS.AddMeal(mealName)
	if err != nil {
		log.Fatal(err)
	}

	insertDiet(d.DB, int64(meal.Id), mealType, source, logForDate)
	fmt.Println("Diet logged successfully")
	return nil
}

func (d *DietServiceImpl) LogDietWithExistingMeal(mealId int64, mealType MealType, source string, logForDate time.Time) error {
	logged, err := d.MealTypeLoggedForToday(mealType)
	if err != nil {
		return err
	}

	if logged {
		return errors.New("meal type already logged for today")
	}

	_, err = meals.MS.GetMeal(int(mealId))
	if err != nil {
		return err
	}

	insertDiet(d.DB, mealId, mealType, source, logForDate)
	fmt.Println("Diet logged successfully")
	return nil
}

func (d *DietServiceImpl) GetDietsForDay(day time.Time) ([]Diet, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.Add(24 * time.Hour)

	var diets []Diet
	err := d.DB.Find(&diets, "timestamp between ? and ?", start.UnixMilli(), end.UnixMilli()).Error
	if err != nil {
		return nil, fmt.Errorf("couldn't fetch diets %v", err)
	}
	return diets, nil
}

func (d *DietServiceImpl) GetDayDietsByOffset(offset int, timestamp int64) ([]DayDiet, error) {
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
		return d.GetBatchedDayDiets(afterUnixMilli, time.Now().UnixMilli())
	}
}

func (d *DietServiceImpl) GetBatchedDayDiets(afterUnixMilli int64, beforeUnixMilli int64) ([]DayDiet, error) {
	diets, err := getBatchedDiets(afterUnixMilli, beforeUnixMilli, d)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	dietMap, err := dietsOfDay(diets)
	if err != nil {
		return nil, fmt.Errorf("couldn't get the diets of the day %v", err)
	}
	return weeklyDiets(dietMap), nil
}

func (d *DietServiceImpl) MealTypeLoggedForToday(mealType MealType) (bool, error) {
	if !ValidMealType(mealType) {
		return false, fmt.Errorf("invalid meal type %v, please use breakfast, lunch, snacks, or dinner", mealType)
	}
	var count int64
	start, end := timeStampRangeForToday()
	if err := d.DB.Model(&Diet{}).Where("meal_type = ? and timestamp between ? and ?", mealType.String(), start, end).Count(&count).Error; err != nil {
		log.Fatal(err)
	}
	return count > 0, nil
}

func (d *DietServiceImpl) EarliestDietTimestamp() int64 {
	var earliest int64
	err := d.DB.Table("diets").Select("min(timestamp)").Scan(&earliest).Error
	if err != nil {
		log.Fatal(err)
	}
	return earliest
}

func ValidMealType(mealType MealType) bool {
	switch mealType {
	case Breakfast, Lunch, Dinner:
		return true
	default:
		return false
	}
}

func GetMealType(mealType string) (MealType, error) {
	switch mealType {
	case "breakfast":
		return Breakfast, nil
	case "lunch":
		return Lunch, nil
	case "dinner":
		return Dinner, nil
	default:
		return 0, fmt.Errorf("invalid meal type %v, please use breakfast, lunch, snacks, or dinner", mealType)
	}
}

func insertDiet(db *gorm.DB, mealId int64, mealType MealType, source string, logForDate time.Time) (Diet, error) {
	diet := Diet{
		MealId:         mealId,
		MealType:       mealType.String(),
		Source:         source,
		Timestamp:      time.Now().UnixMilli(),
		LoggedForDay:   logForDate.Day(),
		LoggedForMonth: int(logForDate.Month()),
		LoggedForYear:  logForDate.Year(),
	}
	if err := db.Create(&diet).Error; err != nil {
		return diet, fmt.Errorf("couldn't log diet %v", err)
	}
	return diet, nil
}

func getBatchedDiets(afterUnixMilli int64, beforeUnixMilli int64, d *DietServiceImpl) ([]Diet, error) {
	var err error
	var diets []Diet

	err = d.DB.Order("timestamp desc").Find(&diets, "timestamp between ? and ?", afterUnixMilli, beforeUnixMilli).Error
	if err != nil {
		return nil, fmt.Errorf("couldn't fetch diets %v", err)
	}
	return diets, nil
}

/*
this function takes a map of day - diets of that day for a week and
returns tabular view of the details of the diets
*/
func weeklyDiets(dietMap map[string][]DietResp) []DayDiet {
	var weekDiets []DayDiet
	for day, diets := range dietMap {

		var (
			d             DayDiet
			totalProtein  int
			totalCarbs    int
			totalFat      int
			totalCalories int
		)

		d.Day = day
		for _, diet := range diets {
			switch diet.MealType {
			case "breakfast":
				d.Breakfast = diet.Meal.Name
			case "lunch":
				d.Lunch = diet.Meal.Name
			case "dinner":
				d.Dinner = diet.Meal.Name
			}

			totalProtein += diet.Meal.GetMealMacro(meals.Protein)
			totalCarbs += diet.Meal.GetMealMacro(meals.Carbs)
			totalFat += diet.Meal.GetMealMacro(meals.Fat)
			totalCalories += diet.Meal.GetTotalCalories()
		}
		d.TotalCalories = totalCalories

		weekDiets = append(weekDiets, d)
	}

	sort.SliceStable(weekDiets, func(i, j int) bool {
		dateI, errI := time.Parse(time.DateOnly, weekDiets[i].Day)
		dateJ, errJ := time.Parse(time.DateOnly, weekDiets[j].Day)
		if errI != nil || errJ != nil {
			log.Fatal("couldn't parse the time")
		}
		return !dateI.Before(dateJ)
	})

	return weekDiets
}

func dietsOfDay(diets []Diet) (map[string][]DietResp, error) {
	dietMap := make(map[string][]DietResp)

	for _, diet := range diets {
		meal, err := meals.MS.GetMeal(int(diet.MealId))
		if err != nil {
			return nil, fmt.Errorf("couldn't get the meal for id: %d %v", diet.MealId, err)
		}

		day := time.UnixMilli(diet.Timestamp).Format("2006-01-02")
		dietMap[day] = append(dietMap[day], DietResp{Day: day, MealType: diet.MealType, Meal: meal})
	}
	return dietMap, nil
}

func timeStampRangeForToday() (int64, int64) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	return startOfDay.UnixMilli(), endOfDay.UnixMilli()
}
