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
