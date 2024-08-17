`cmdiet` is a command line tool to log your diets with swag! It's target audience are programmers who never want to leave their terminals (honestly why would you?). It is built on the [bubbletea](https://github.com/charmbracelet/bubbletea) package and that's why it's, you guessed it, delightful to use.

- Log your breakfast, lunch and dinner with minimum inputs
- View your last week's diets, including calories and macro info. Or provide a `batch` option to see last `n` days of diets.
- Use the `left` and `right` keys to go into the past or future if you are at a batch in the past
- Own your catalogue of meals. View and edit your meals in a beautiful UI. You can also filter through them with a fuzzy search.
- Edit the calorie and macro info for a meal anytime. The updated info will be instantly reflected in your logs.
- Analyze your calorie and macro changes in a nice-looking plot. Go as far as you want in the past. Switch between braille and histogram view with a key.

## Install

### Via package managers

```bash
# macOS or Linux
brew install cmdiet

# Windows
choco install cmdiet
winget install cmdiet
```

### Using Go

```
go get
```

### Log a meal

```
cmdiet log {breakfast | lunch | dinner | misc}
```

Enter the meal name or select from your meal catalogue. If you enter a meal, it will be added to your meal catalogue with empty values for protein, carbohydrates and fats. You can edit these value from the `list` command.

**Options:**

`misc` : Log a miscellaneous meal. Anything that you eat outside of your regular meals.

### Edit your meal catalogue

```
cmdiet list {meals}
```

List all your meals. You can edit the calorie and macro info for a meal from here. Just enter on a meal to edit it.
You can also filter through your meals with a fuzzy search.

### View your logs

```
cmdiet view [today] [-o, --offset] [week]
```

View your logs for today or an offset from today. You can also view logs for the last week with the `week` option. Or provide a `batch` option to see last `n` days of diets.

Use left and right keys to go into the past or future if you are at a batch in the past.

Enter on a log to view the summary of that day.

**Options:**

`today`: View today's diet logs, macros breakdown and calorie breakdown by macros

`-o, --offset`: View logs in a tabular view for an offset numbers of days in past starting from today

`week`: View logs for the last week. Default option if no option is provided.

### Analyze your logs

```
cmdiet eval {calories|protein|carbohydrates|fats}
```

Analyze your calorie and macro changes in a nice-looking plot. Use the left and right keys to go as far as you want in the past. Switch between braille and histogram view with a key.
