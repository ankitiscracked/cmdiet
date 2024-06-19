package database

var CreateDietTableSql = `create table if not exists diet (
  id integer primary key autoincrement,
  meal_id integer,
  meal_type text,
  source text,
  timestamp integer,
  foreign key (meal_id) references meal(id)
);`

var CreateMealTableSql = `create table if not exists meals (
  id integer primary key autoincrement,
  name text,
  calories integer,
  protein integer,
  carbs integer,
  fats integer,
  timestamp integer
);`
